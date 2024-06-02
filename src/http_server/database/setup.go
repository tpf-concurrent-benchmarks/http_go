package database

import (
	"database/sql"
	"fmt"
	"os"
	"github.com/gin-gonic/gin"
	"time"
	_ "github.com/lib/pq"
)

// part of the code was taken from https://nerocui.com/2019/08/04/how-to-use-sql-database-in-golang/
func InitializeDatabase() *sql.DB {
	host, port, user, password, name := getConfig()
	var err error
	psqlInfo := fmt.Sprintf(`
		host=%s port=%s user=%s 
		password=%s dbname=%s sslmode=disable`,
		host, port, user, password, name)
	var db *sql.DB
	for { // Retry until connection is established
		db, err = sql.Open("postgres", psqlInfo)
		if err != nil {
			fmt.Printf("Error opening database: %v. Retrying in 3 seconds...\n", err)
			time.Sleep(3 * time.Second)
			continue
		}
		
		err = db.Ping() // Check if connection is established
		if err != nil {
			fmt.Printf("Error pinging database: %v. Retrying in 3 seconds...\n", err)
			time.Sleep(3 * time.Second)
			continue
		}

		fmt.Println("Connected to PostgreSQL")
		break
	}
	
	fmt.Println("Connected to PostgreSQL")
	err = createTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

func DatabaseMiddleware(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db) // Attach db to context
		c.Next()
	}
}

// gets the database connection information from the environment variables
func getConfig() (string, string, string, string, string) {
	host := os.Getenv("DBHOST")
	port := os.Getenv("DBPORT")
	user := os.Getenv("DBUSER")
	password := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")
	return host, port, user, password, name
}

// gets the db from the gin context
func getDB(c *gin.Context) *sql.DB {
	return c.MustGet("db").(*sql.DB)
}

// this extension is used to create unique user ids
func activateExtension(db *sql.DB) error {
	sqlStatement := `CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`
	_, err := db.Exec(sqlStatement)
	return err
}