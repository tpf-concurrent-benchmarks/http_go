package http_server

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"github.com/fergusstrange/embedded-postgres"
	"github.com/gin-gonic/gin"
	models "http_go/http_server/models"
)

func getConfig() (string, string, string, string, string) {
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")
	return host, port, user, password, name
}

func createTables(db *sql.DB) error {
	err := activateExtension(db)
	if err != nil {
		return err
	}
	err = createUserTable(db)
	if err != nil {
		return err
	}
	err = createPollTable(db)
	if err != nil {
		return err
	}
	err = createPollOptionsTable(db)
	if err != nil {
		return err
	}
	err = createVotesTable(db)
	if err != nil {
		return err
	}
	return nil
}

// part of the code was taken from https://nerocui.com/2019/08/04/how-to-use-sql-database-in-golang/
func InitializeDatabase() *sql.DB {
	host, port, user, password, name := getConfig()
	var err error
	psqlInfo := fmt.Sprintf(`
		host=%s port=%s user=%s 
		password=%s dbname=%s sslmode=disable`,
		host, port, user, password, name)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to PostgreSQL")
	err = createTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func StartDatabase() *embeddedpostgres.EmbeddedPostgres {
	postgres := embeddedpostgres.NewDatabase()
	err := postgres.Start()
	if err != nil {
		panic(err)
	}
	return postgres
}


func CloseDatabase(postgres *embeddedpostgres.EmbeddedPostgres) {
	err := postgres.Stop()
	if err != nil {
		panic(err)
	}
}

func DatabaseMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db) // Attach db to context
		c.Next()
	}
}

func GetDB(c *gin.Context) *sql.DB {
	return c.MustGet("db").(*sql.DB)
}

func activateExtension(db *sql.DB) error {
	sqlStatement := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	_, err := db.Exec(sqlStatement)
	return err
}

func createUserTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE users (
		id TEXT PRIMARY KEY,
		username CHAR(30) NOT NULL,
		password CHAR(64) NOT NULL
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func createPollTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE polls (
		poll_id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
		creator_id UUID NOT NULL,
		poll_topic TEXT NOT NULL
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func createPollOptionsTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE poll_options (
		poll_id UUID,
		option_num INT,
		option_text TEXT NOT NULL,
		PRIMARY KEY (poll_id, option_num)
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func createVotesTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE votes (
		poll_id UUID,
		user_id UUID,
		option_num INT NOT NULL,
		PRIMARY KEY (poll_id, user_id)
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func insertUser(c *gin.Context, username string, password string) error {
	//TODO: check if user already exists
	db := GetDB(c)
	sqlStatement := `INSERT INTO users (id, username, password) VALUES (uuid_generate_v4(), $1, $2)`
	_, err := db.Exec(sqlStatement, username, password)
	return err
}

func getUserPassword(c *gin.Context, username string) (string, error) {
	db := GetDB(c)
	sqlStatement := `SELECT password FROM users WHERE username=$1`
	var password string
	err := db.QueryRow(sqlStatement, username).Scan(&password)
	return password, err
}

func getUserID(c *gin.Context, username string) (string, error) {
	db := GetDB(c)
	sqlStatement := `SELECT id FROM users WHERE username=$1`
	var id string
	err := db.QueryRow(sqlStatement, username).Scan(&id)
	return id, err
}

func getUser(c *gin.Context, username string) (models.UserData, error) {
	db := GetDB(c)
	sqlStatement := `SELECT id, username, password FROM users WHERE username=$1`
	var user models.UserData
	err := db.QueryRow(sqlStatement, username).Scan(&user.ID, &user.Username, &user.HashedPassword)
	return user, err
}

func insertPoll(c *gin.Context, creatorID string, poll models.Poll) (string, error) {
	db := GetDB(c)
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

func getPoll(c *gin.Context, pollID string) (models.PollWithVotes, error) {
	db := GetDB(c)
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

func insertVote(c *gin.Context, vote models.Vote, userID string) error {
	db := GetDB(c)
	sqlStatement := `INSERT INTO votes (poll_id, user_id, option_num) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, vote.PollID, userID, vote.Option)
	return err
}
