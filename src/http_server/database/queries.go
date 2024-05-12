package database

import (
	"github.com/gin-gonic/gin"
	models "http_go/http_server/models"
	"fmt"
	"database/sql"
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
	//check poll exists
	sqlStatement := `SELECT COUNT(*) FROM polls WHERE poll_id=$1`
	var count int
	err := db.QueryRow(sqlStatement, vote.PollID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("poll does not exist")
	}
	//TODO: check if this vote exists
	sqlStatement = `INSERT INTO votes (poll_id, user_id, option_num) VALUES ($1, $2, $3)`
	_, err = db.Exec(sqlStatement, vote.PollID, userID, vote.Option)
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
	printTableContents(db, "votes")
	sqlStatement = `SELECT po.option_text, COUNT(v.poll_id) AS vote_count
					FROM poll_options po
						LEFT JOIN votes v ON po.poll_id = v.poll_id 
						AND po.option_num = v.option_num
					WHERE po.poll_id = $1
					GROUP BY po.option_text`
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

func GetPolls(c *gin.Context) ([]models.PollMeta, error) {
	db := getDB(c)
	printTableContents(db, "polls")
	sqlStatement := `SELECT poll_id, poll_topic FROM polls`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	fmt.Println(rows)
	defer rows.Close()
	var polls []models.PollMeta
	for rows.Next() {
		fmt.Println("in loop")
		var poll models.PollMeta
		err = rows.Scan(&poll.ID, &poll.Title)
		if err != nil {
			continue
		}
		polls = append(polls, poll)
	}
	fmt.Println(polls)
	return polls, nil
}

func printTableContents(db *sql.DB, tableName string) error {
	// Construct the SQL query to select all rows from the specified table
	query := fmt.Sprintf("SELECT * FROM %s", tableName)

	// Execute the SQL query using db.Query
	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Get column names from the result set
	columnNames, err := rows.Columns()
	if err != nil {
		return err
	}

	// Print column names
	fmt.Println("Table:", tableName)
	fmt.Println("Columns:", columnNames)

	// Iterate over the rows returned by the query
	for rows.Next() {
		// Slice to hold column values for this row
		columnValues := make([]interface{}, len(columnNames))
		columnPointers := make([]interface{}, len(columnNames))

		// Map column values to column pointers for Scan method
		for i := range columnValues {
			columnPointers[i] = &columnValues[i]
		}

		// Scan row data into columnValues slice
		if err := rows.Scan(columnPointers...); err != nil {
			return err
		}

		// Print row data
		for i, value := range columnValues {
			fmt.Printf("%s: %v\t", columnNames[i], value)
		}
		fmt.Println() // Newline after each row
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}