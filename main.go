package main

import (
	"energytoken/controllers"
	"energytoken/middlewares"
	"energytoken/models"

	// "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	models.ConnectDatabase()

	r := gin.Default()

	public := r.Group("/api", CORSMiddleware())

	public.OPTIONS("/migrate", CORSMiddleware())
	public.GET("/migrate", controllers.Migrate)

	public.OPTIONS("/register", CORSMiddleware())
	public.POST("/register", controllers.Register)

	public.OPTIONS("/login", CORSMiddleware())
	public.POST("/login", controllers.Login)

	protectedUsers := r.Group("/api/user", CORSMiddleware())
	protectedUsers.Use(middlewares.JWTAuthMiddleware())

	protectedUsers.OPTIONS("/transfer", CORSMiddleware())
	protectedUsers.POST("/transfer", controllers.Transfer)

	protectedUsers.OPTIONS("/balance", CORSMiddleware())
	protectedUsers.GET("/balance", controllers.Balance)

	// should be admin
	protected := r.Group("/api/admin", CORSMiddleware())
	protected.Use(middlewares.JWTAuthMiddleware())
	protected.OPTIONS("/user", CORSMiddleware())
	protected.POST("/user", controllers.CurrentUser)

	r.Run(":8080")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
