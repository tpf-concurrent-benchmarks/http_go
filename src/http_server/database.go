package http_server

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"github.com/fergusstrange/embedded-postgres"
	"github.com/google/uuid"
)

func getConfig() (string, string, string, string, string) {
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")
	return host, port, user, password, name
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
	err = createUserTable(db)
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

func createUserTable(db *sql.DB) error {
	sqlStatement := `
	CREATE TABLE users (
		id TEXT PRIMARY KEY,
		username TEXT NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := db.Exec(sqlStatement)
	return err
}

func insertUser(c *gin.Context,  id UUID, username string, password string) error {
	db := GetDB(c)
	sqlStatement := `INSERT INTO users (id, username, password) VALUES ($1, $2, $3)`
	_, err := db.Exec(sqlStatement, id.String(), username, password)
	return err
}