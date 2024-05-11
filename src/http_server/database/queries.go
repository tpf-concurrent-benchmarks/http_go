package database

import (
	"github.com/gin-gonic/gin"
	models "http_go/http_server/models"
)

func InsertUser(c *gin.Context, username string, password string) error {
	//TODO: check if user already exists
	db := getDB(c)
	sqlStatement := `INSERT INTO users (id, username, password) VALUES (uuid_generate_v4(), $1, $2)`
	_, err := db.Exec(sqlStatement, username, password)
	return err
}

func GetUserPassword(c *gin.Context, username string) (string, error) {
	db := getDB(c)
	sqlStatement := `SELECT password FROM users WHERE username=$1`
	var password string
	err := db.QueryRow(sqlStatement, username).Scan(&password)
	return password, err
}

func GetUserID(c *gin.Context, username string) (string, error) {
	db := getDB(c)
	sqlStatement := `SELECT id FROM users WHERE username=$1`
	var id string
	err := db.QueryRow(sqlStatement, username).Scan(&id)
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

func GetPoll(c *gin.Context, pollID string) (models.PollWithVotes, error) {
	db := getDB(c)
	sqlStatement := `SELECT poll_topic FROM polls WHERE poll_id=$1`
	var poll models.Poll
	err := db.QueryRow(sqlStatement, pollID).Scan(&poll.Title)
	if err != nil {
		return models.PollWithVotes{}, err
	}
	sqlStatement = `SELECT option_text FROM poll_options WHERE poll_id=$1`
	rows, err := db.Query(sqlStatement, pollID)
	if err != nil {
		return models.PollWithVotes{}, err
	}
	defer rows.Close()
	var options []models.Option
	for rows.Next() {
		var option models.Option
		err = rows.Scan(&option.Name)
		if err != nil {
			return models.PollWithVotes{}, err
		}
		options = append(options, option)
	}
	return models.PollWithVotes{Title: poll.Title, Options: options}, nil
}

func InsertVote(c *gin.Context, vote models.Vote, userID string) error {
	db := getDB(c)
	sqlStatement := `INSERT INTO votes (poll_id, user_id, option_num) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, vote.PollID, userID, vote.Option)
	return err
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
	sqlStatement = `SELECT option_text, COUNT(option_num) 
					FROM poll_options LEFT JOIN votes 
						ON poll_options.poll_id=votes.poll_id 
						AND poll_options.option_num=votes.option_num 
					WHERE poll_options.poll_id=$1 
					GROUP BY option_text`
	rows, err := db.Query(sqlStatement, pollID)
	if err != nil {
		return models.PollWithVotes{}, err
	}
	defer rows.Close()
	var options []models.Option
	for rows.Next() {
		var option models.Option
		err = rows.Scan(&option.Name, &option.Votes)
		if err != nil {
			continue
		}
		options = append(options, option)
	}
	return models.PollWithVotes{Title: poll.Title, Options: options}, nil
}