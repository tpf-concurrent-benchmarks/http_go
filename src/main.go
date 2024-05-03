package main

import (
   "github.com/gin-gonic/gin"
   swaggerfiles "github.com/swaggo/files"
   ginSwagger "github.com/swaggo/gin-swagger"
   "net/http"
   docs "http_go/docs"
   "fmt"
   "github.com/golang-jwt/jwt/v5"
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
// @Param username body string true "Username to add"
// @Success 200 {string} string "User added successfully"
// @Failure 400 {string} string "Invalid request payload"
// @Router /users [post]
func create_user(c *gin.Context) {
	var username string

	// Bind JSON request body to username variable
	if err := c.ShouldBindJSON(&username); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// Check if username already exists
	if _, exists := db[username]; exists {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}

	// Add username to the database
	db[username] = "true"

	c.JSON(200, gin.H{"message": "User added successfully"})
}

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	// Ping test
	v1.GET("/ping", Helloworld)

	
	v1.POST("/users", create_user)
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
