package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "http_go/docs"
	server "http_go/http_server"
	"database/sql"
	db "http_go/http_server/database"
)

// sets up the routes of the api
func setupRouter(db_controller *sql.DB) *gin.Engine {
	jwtManager := server.NewJWTManager()

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api"
	v1 := r.Group("/api")
	v1.Use(db.DatabaseMiddleware(db_controller))
	
	v1.POST("/users", func(c *gin.Context) {
		server.CreateUser(jwtManager, c)
	})
	v1.POST("/login", func(c *gin.Context) {
		server.Login(jwtManager, c)
	})
	v1.POST("/polls", func(c *gin.Context) {
		server.CreatePoll(jwtManager, c)
	})
	v1.GET("/polls/:id", server.GetPoll)
	v1.GET("/polls", server.GetPolls)
	v1.POST("polls/:id/vote", func(c *gin.Context) {
		server.Vote(jwtManager, c)
	})
	v1.DELETE("/polls/:id", func(c *gin.Context) {
		server.DeletePoll(jwtManager, c)
	})

	return r
}

func main() {
	server.LoadEnv()

	db_controller := db.InitializeDatabase()
	r := setupRouter(db_controller)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Listen and Server in 0.0.0.0:8765
	r.Run(":8765")
}
