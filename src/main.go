package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "http_go/docs"
	server "http_go/http_server"
)

func setupRouter() *gin.Engine {
	jwtManager := server.NewJWTManager()

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")

	
	v1.POST("/users", server.CreateUser)
	v1.POST("/login", func(c *gin.Context) {
		server.Login(jwtManager, c)
	})
	v1.POST("/poll", func(c *gin.Context) {
		server.CreatePoll(jwtManager, c)
	})
	v1.GET("/poll/:id", server.GetPoll)
	v1.GET("/users/:name", server.UserExists)

	return r
}

func main() {
	r := setupRouter()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
