package http_server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	models "http_go/http_server/models"
	"time"
)

var db = make(map[string]string)

// @BasePath /api/v1

// @Router /users/:name [get]
// @Param  name query string true "name"
func UserExists(c *gin.Context)  {
	user := c.Request.URL.Query().Get("name")
	value, ok := db[user]
	if ok {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "value": value})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "value": "not found"})
	}
}

// addUser adds a username to the database
// @Summary Add a new user
// @Description Add a new user to the database
// @Accept json
// @Produce json
// @Param userInDB body models.UserInDB true "Username and hashed password"
// @Success 200 {string} string "User added successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var user models.UserInDB

	// Bind JSON request body to username variable
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if username already exists
	if _, exists := db[user.Username]; exists {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	// Add username to the database
	db[user.Username] = user.HashedPassword

	c.JSON(200, gin.H{"message": "User added successfully"})
}

// @Router /login [post]
// @Param userInDB body models.UserInDB true "Username and hashed password"
// @Success 200 {string} string "access_token"
// @Failure 400 {string} string "Invalid request body"
func Login(jwtManager *JWTManager, c *gin.Context) {
	var user models.UserInDB

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userPassword, ok := db[user.Username]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if userPassword != user.HashedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenData := map[string]interface{}{
		"sub": user.Username,
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
}