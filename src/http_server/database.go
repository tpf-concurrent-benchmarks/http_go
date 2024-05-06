package http_server

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
	"os"
	"github.com/fergusstrange/embedded-postgres"

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
	return db
}

func StartDatabase() *embeddedpostgres.EmbeddedPostgres {
	postgres := embeddedpostgres.NewDatabase()
	err := postgres.Start()
	if err != nil {
		panic(err)
	}
	// fmt.Println("started database")
	return postgres
}


func CloseDatabase(postgres *embeddedpostgres.EmbeddedPostgres) {
	err := postgres.Stop()
	if err != nil {
		panic(err)
	}
}