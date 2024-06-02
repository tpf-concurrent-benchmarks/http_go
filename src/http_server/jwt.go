package http_server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
	"strings"
)

type JWTManager struct {
	SecretKey                string
	AccessTokenExpireMinutes int
	Algorithm                jwt.SigningMethod
}

// initializes the configuration for the jwt tokens
func NewJWTManager() *JWTManager {
	if os.Getenv("ALGORITHM") == "HS256" {
		return &JWTManager{
			SecretKey:                os.Getenv("SECRET_KEY"),
			AccessTokenExpireMinutes: getEnvInt("ACCESS_TOKEN_EXPIRE_MINUTES", 30),
			Algorithm:                jwt.SigningMethodHS256,
		}
	} else {
		return &JWTManager{
			SecretKey:                os.Getenv("SECRET_KEY"),
			AccessTokenExpireMinutes: getEnvInt("ACCESS_TOKEN_EXPIRE_MINUTES", 30),
			Algorithm:                jwt.SigningMethodHS512,
		}
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

// decodes the jwt token, it checks for a valid encoding
// and if the token has expired
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
    if int64(claims["exp"].(float64)) < time.Now().Unix() {
        return nil, fmt.Errorf("token has expired")
    }
	return claims, nil
}

// processes the jwt token from the header Authorization
// the format of the header should be "Bearer <token>"
// it also decodes the token and returns the claims
func processToken(jwtManager *JWTManager, c *gin.Context) (jwt.MapClaims, error) {
	// Check if the "Authorization" header exists in the request
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return nil, fmt.Errorf("authorization header is missing")
	}

	// Split the header value to get the token
	authValues := strings.Split(authHeader, " ")
	if len(authValues) != 2 || authValues[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
		return nil, fmt.Errorf("invalid authorization header format")
	}

	accessToken := authValues[1]

	claims, err := jwtManager.DecodeJWTToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return nil, err
	}

	return claims, nil
}