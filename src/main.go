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

func setupRouter(db_controller *sql.DB) *gin.Engine {
	jwtManager := server.NewJWTManager()

	r := gin.Default()
	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	v1.Use(db.DatabaseMiddleware(db_controller))
	
	v1.POST("/users", server.CreateUser)
	v1.POST("/login", func(c *gin.Context) {
		server.Login(jwtManager, c)
	})
	v1.POST("/poll", func(c *gin.Context) {
		server.CreatePoll(jwtManager, c)
	})
	v1.GET("/poll/:id", server.GetPoll)
	v1.GET("/polls", server.GetPolls)
	v1.POST("poll/:id/vote", func(c *gin.Context) {
		server.Vote(jwtManager, c)
	})
	v1.DELETE("/polls/:id", func(c *gin.Context) {
		server.DeletePoll(jwtManager, c)
	})

	return r
}

func main() {
	server.LoadEnv()
	postgres := db.StartDatabase()
	defer db.CloseDatabase(postgres)
	server.SigHandler(postgres)
	db_controller := db.InitializeDatabase()
	r := setupRouter(db_controller)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Listen and Server in 0.0.0.0:8765
	r.Run(":8765")
}
