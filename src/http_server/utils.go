package http_server

import (
	"os"
	"strconv"
	"github.com/joho/godotenv"
	"log"
	"crypto/sha256"
	"encoding/hex"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// fetch an int from an environment variable
// returns the default value if the key is not found or the value is not an int
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
