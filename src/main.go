package main

import (
   "github.com/gin-gonic/gin"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
   docs "http_go/docs"
//    "fmt"
	jwt "http_go/http_server"
	models "http_go/http_server/models"
	"time"
)

// @BasePath /api/v1

// PingExample godoc
// @Summary ping example
// @Schemes
// @Description do ping
// @Tags example
// @Accept json
// @Produce json
// @Success 200 {string} Helloworld
// @Router /ping [get]
func Helloworld(g *gin.Context)  {
   g.JSON(http.StatusOK,"helloworld")
}

// @Router /users/:name [get]
// @Param  name query string true "name"
func user_exists(c *gin.Context)  {
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
func create_user(c *gin.Context) {
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
func login(jwtManager *jwt.JWTManager, c *gin.Context) {
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


var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	jwtManager := jwt.NewJWTManager()

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	// Ping test
	v1.GET("/ping", Helloworld)

	
	v1.POST("/users", create_user)
	v1.POST("/login", func(c *gin.Context) {
		login(jwtManager, c)
	})
	v1.GET("/users/:name", user_exists)

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := v1.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	r := setupRouter()

	// {
	// 	eg := v1.Group("/example")
	// 	{
	// 		eg.GET("/helloworld",Helloworld)
	// 	}
	// }
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
