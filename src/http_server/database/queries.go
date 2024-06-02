package database

import (
	"github.com/gin-gonic/gin"
	models "http_go/http_server/models"
	"fmt"
	"database/sql"
)

func InsertUser(c *gin.Context, username string, password string) (string, error) {
	//TODO: check if user already exists
	db := getDB(c)
	sqlStatement := `INSERT INTO users (id, username, password) VALUES (uuid_generate_v4(), $1, $2) RETURNING id`
	var id string
	err := db.QueryRow(sqlStatement, username, password).Scan(&id)
	return id, err
}

func GetUser(c *gin.Context, username string) (models.UserData, error) {
	db := getDB(c)
	sqlStatement := `SELECT id, username, password FROM users WHERE username=$1`
	var user models.UserData
	err := db.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &user.HashedPassword)
	return user, err
}

func InsertPoll(c *gin.Context, creatorID string, poll models.Poll) (string, error) {
	db := getDB(c)
	sqlStatement := `INSERT INTO polls (poll_id, creator_id, poll_topic)
					VALUES (uuid_generate_v4(), $1, $2)
					RETURNING poll_id`
	var pollID string
	err := db.QueryRow(sqlStatement, creatorID, poll.Title).Scan(&pollID)
	if err != nil {
		return "", err
	}
	for i, option := range poll.Options {
		sqlStatement = `INSERT INTO poll_options (poll_id, option_num, option_text) VALUES ($1, $2, $3)`
		_, err = db.Exec(sqlStatement, pollID, i, option)
	}
	return pollID, err
}

func InsertVote(c *gin.Context, vote models.Vote, userID string) error {
	db := getDB(c)

	// Check if poll exists
	sqlStatement := `SELECT COUNT(*) FROM polls WHERE poll_id = $1`
	var count int
	err := db.QueryRow(sqlStatement, vote.PollID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("poll does not exist")
	}

	// Check if the vote already exists
	var existingOptionNum int
	sqlStatement = `SELECT option_num FROM votes WHERE poll_id = $1 AND user_id = $2`
	err = db.QueryRow(sqlStatement, vote.PollID, userID).Scan(&existingOptionNum)

	switch {
	case err == sql.ErrNoRows:
		// Vote does not exist, insert a new row
		sqlStatement = `INSERT INTO votes (poll_id, user_id, option_num)
						VALUES ($1, $2, $3)`
		_, err = db.Exec(sqlStatement, vote.PollID, userID, vote.Option)
		return err

	case err != nil:
		// Other error occurred
		return err

	default:
		// Vote exists, check the option_num
		if existingOptionNum != vote.Option {
			// Option number is different, update the option_num
			sqlStatement = `UPDATE votes SET option_num = $1 WHERE poll_id = $2 AND user_id = $3`
			_, err = db.Exec(sqlStatement, vote.Option, vote.PollID, userID)
			return err
		} else {
			// Option number is the same, delete the existing vote
			sqlStatement = `DELETE FROM votes WHERE poll_id = $1 AND user_id = $2`
			_, err = db.Exec(sqlStatement, vote.PollID, userID)
			return err
		}
	}
}

func GetPollWithVotes(c *gin.Context, pollID string) (models.PollWithVotes, error) {
	db := getDB(c)
	sqlStatement := `SELECT poll_topic 
					 FROM polls 
					 WHERE poll_id=$1`
	var poll models.Poll
	err := db.QueryRow(sqlStatement, pollID).Scan(&poll.Title)
	if err != nil {
		return models.PollWithVotes{}, err
	}
	sqlStatement = `SELECT po.option_text, COUNT(v.poll_id) AS vote_count
					FROM poll_options po
						LEFT JOIN votes v ON po.poll_id = v.poll_id 
						AND po.option_num = v.option_num
					WHERE po.poll_id = $1
					GROUP BY po.option_text, po.option_num
					ORDER BY po.option_num ASC` //this may cause a problem if two options have the same text
	rows, err := db.Query(sqlStatement, pollID)
	if err != nil {
		return models.PollWithVotes{}, err
	}
	defer rows.Close()
	var options []models.Option
	for rows.Next() {
		var option models.Option
		option.Votes = 0
		err = rows.Scan(&option.Name, &option.Votes)
		if err != nil {
			continue
		}
		options = append(options, option)
	}
	return models.PollWithVotes{Title: poll.Title, Options: options}, nil
}

func GetPolls(c *gin.Context) ([]models.PollMeta, error) {
	db := getDB(c)
	sqlStatement := `SELECT poll_id, poll_topic FROM polls`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var polls []models.PollMeta
	for rows.Next() {
		var poll models.PollMeta
		err = rows.Scan(&poll.ID, &poll.Title)
		if err != nil {
			continue
		}
		polls = append(polls, poll)
	}
	return polls, nil
}

func GetPollCreator(c *gin.Context, pollID string) (string, error) {
	db := getDB(c)
	sqlStatement := `SELECT creator_id FROM polls WHERE poll_id=$1`
	var creatorID string
	err := db.QueryRow(sqlStatement, pollID).Scan(&creatorID)
	return creatorID, err
}

func DeletePoll(c *gin.Context, pollID string) error {
	db := getDB(c)
	sqlStatement := `DELETE FROM polls WHERE poll_id = $1`
	_, err := db.Exec(sqlStatement, pollID)
	combinedErr := ""
	if err != nil {
		combinedErr += "delete error in polls: " + err.Error() + "\n"

	}
	sqlStatement = `DELETE FROM poll_options WHERE poll_id = $1`
	_, err = db.Exec(sqlStatement, pollID)
	if err != nil {
		combinedErr += "delete error in poll_options: " + err.Error() + "\n"
	}
	sqlStatement = `DELETE FROM votes WHERE poll_id = $1`
	_, err = db.Exec(sqlStatement, pollID)
	if err != nil {
		combinedErr += "delete error in votes: " + err.Error() + "\n"
	}
	if combinedErr != "" {
		return fmt.Errorf(combinedErr)
	}
	return nil
}