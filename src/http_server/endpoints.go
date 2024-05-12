package http_server

import (
	"github.com/gin-gonic/gin"
	"net/http"
	models "http_go/http_server/models"
	"time"
	db "http_go/http_server/database"
)

// @BasePath /api/v1

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

	// Add username to the database
	err := db.InsertUser(c, user.Username, user.HashedPassword)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add user"})
		return
	}

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

	userData, err := db.GetUser(c, user.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	if userData.HashedPassword != user.HashedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	tokenData := map[string]interface{}{
		"sub": user.Username,
		"id":  userData.ID,
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

// @Router /poll [post]
// @Param token header models.Token true "Bearer token"
// @Param poll body models.Poll true "Poll object"
// @Success 200 {string} string "Poll created successfully"
// @Failure 400 {string} string "Invalid request payload"
func CreatePoll(jwtManager *JWTManager, c *gin.Context) {
	
	claims, err := processToken(jwtManager, c)
	if err != nil { return }
	user_id := claims["id"].(string)
	var poll models.Poll
	if err := c.ShouldBindJSON(&poll); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	ID, err := db.InsertPoll(c, user_id, poll)

	c.JSON(200, gin.H{"message": "Poll created successfully", "id": ID})
}

// @Router /poll/{id} [get]
// @Param id path string true "Poll ID"
// @Success 200 {string} string "PollWithVotes object"
// @Failure 404 {string} string "Poll not found"
func GetPoll(c *gin.Context) {
	ID := c.Param("id")
	poll, err := db.GetPollWithVotes(c, ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error getting poll", "message": err.Error()})
		return
	}
	c.JSON(200, poll)
}

// @Router /polls [get]
// @Success 200 {string} string "Polls object"
// @Failure 404 {string} string "Polls not found"
func GetPolls(c *gin.Context) {
	polls, err := db.GetPolls(c)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error getting polls", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"polls": polls})
}

// @Router /poll/{id}/vote [post]
// @Param token header models.Token true "Bearer token"
// @Param vote body models.Vote true "Vote object"
// @Success 200 {string} string "Voted successfully"
// @Failure 400 {string} string "Invalid request payload"
func Vote(jwtManager *JWTManager, c *gin.Context) {
	var vote models.Vote
	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	claims, err := processToken(jwtManager, c)
	if err != nil { return }
	if claims["sub"].(string) != vote.Username {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = db.InsertVote(c, vote, claims["id"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to vote", "message": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Voted successfully"})
}

// @Router /polls/{id} [delete]
// @Param token header models.Token true "Bearer token"
// @Param id path string true "Poll ID"
// @Success 200 {string} string "Poll deleted successfully"
// @Failure 404 {string} string "Poll not found"
func DeletePoll(jwtManager *JWTManager, c *gin.Context) {
	claims, err := processToken(jwtManager, c)
	if err != nil { return }
	user_id := claims["id"].(string)

	poll_id := c.Param("id")
	creator_id, err := db.GetPollCreator(c, poll_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "error getting poll", "message": err.Error()})
		return
	}
	if creator_id != user_id {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, user not recognized as creator of poll"})
		return
	} 
	err = db.DeletePoll(c, poll_id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Error deleting poll", "message": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Poll deleted successfully"})
}