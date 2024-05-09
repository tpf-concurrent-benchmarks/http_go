package http_server

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

type JWTManager struct {
	SecretKey                string
	AccessTokenExpireMinutes int
	Algorithm                jwt.SigningMethod
}

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

func processToken(jwtManager *JWTManager, c *gin.Context) (jwt.MapClaims, error) {
	// Check if the "Access_token" header exists in the request
	if accessTokenValues, ok := c.Request.Header["Access_token"]; !ok || len(accessTokenValues) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token is missing"})
		return nil, fmt.Errorf("access token is missing")
	}
	if tokenTypeValues, ok := c.Request.Header["Access_token"]; !ok || len(tokenTypeValues) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token type is missing"})
		return nil, fmt.Errorf("token type is missing")
	}
	access_token := c.Request.Header["Access_token"][0]
	token_type := c.Request.Header["Token_type"][0]
	if token_type != "bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Un supported token type"})
		return nil, fmt.Errorf("unsupported token type")
	}

	claims, err := jwtManager.DecodeJWTToken(access_token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return nil, err
	}
	return claims, nil
}