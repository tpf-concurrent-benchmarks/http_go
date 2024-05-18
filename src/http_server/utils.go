package http_server

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
	"log"
	"crypto/sha256"
	"encoding/hex"
	db "http_go/http_server/database"
	"fmt"
	"github.com/fergusstrange/embedded-postgres"
	"os/signal"
	"syscall"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func getEnvInt(key string, defaultValue int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	hashedBytes := hash.Sum(nil)
	return hex.EncodeToString(hashedBytes)
}

func cleanup(postgres *embeddedpostgres.EmbeddedPostgres) {
	db.CloseDatabase(postgres)
    fmt.Println("cleanup")
}

func SigHandler(postgres *embeddedpostgres.EmbeddedPostgres) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cleanup(postgres)
		os.Exit(1)
	}()
}
