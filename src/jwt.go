package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gin-gonic/gin"
	"./models"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

type JWTManager struct {
	SecretKey                string
	AccessTokenExpireMinutes int
	Algorithm                jwt.SigningMethod
}

func NewJWTManager() *JWTManager {
	return &JWTManager{
		SecretKey:                os.Getenv("SECRET_KEY"),
		AccessTokenExpireMinutes: getEnvInt("ACCESS_TOKEN_EXPIRE_MINUTES", 30),
		Algorithm:                jwt.SigningMethodHS256,
	}
}

func (jm *JWTManager) CreateJWTToken(data map[string]interface{}, expiresDelta time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Add(expiresDelta).Unix(),
	}
	for key, value := range data {
		claims[key] = value
	}
	token := jwt.NewWithClaims(jm.Algorithm, claims)
	return token.SignedString([]byte(jm.SecretKey))
}

// DecodeJWTToken decodes and verifies a JWT token, returning the claims if valid
func (jm *JWTManager) DecodeJWTToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jm.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}
	return claims, nil
}

// getEnvInt retrieves an integer environment variable or returns a default value
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
func main() {
	jwtManager := NewJWTManager()

	router := gin.Default()

	// Login endpoint to generate JWT token
	router.POST("/api/login", func(c *gin.Context) {
		var credentials struct {
			Username string `json:"username"`
			// Add other fields as needed (e.g., password)
		}
		if err := c.ShouldBindJSON(&credentials); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		// Validate user credentials (e.g., authenticate against database)
		// For demonstration, assuming authentication is successful
		tokenData := map[string]interface{}{
			"sub": credentials.Username,
			// Add other claims as needed
		}

		token, err := jwtManager.CreateJWTToken(tokenData, time.Minute*time.Duration(jwtManager.AccessTokenExpireMinutes))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token": token,
			"token_type":   "bearer",
		})
	})

	// Protected endpoint to get user from token
	router.GET("/api/user", func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")[len("Bearer "):]
		claims, err := jwtManager.DecodeJWTToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			return
		}

		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().UTC().After(expirationTime) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			return
		}

		username := claims["sub"].(string)
		c.JSON(http.StatusOK, gin.H{"username": username})
	})

	// Run the server
	router.Run(":8080")
}