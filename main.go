package main

import (
	// "time"

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

	public.OPTIONS("/registercheck", CORSMiddleware())
	public.POST("/registercheck", controllers.ValidatePhoneAndID)

	public.OPTIONS("/register", CORSMiddleware())
	public.POST("/register", controllers.Register)

	public.OPTIONS("/login", CORSMiddleware())
	public.POST("/login", controllers.Login)

	public.OPTIONS("/verifylogin", CORSMiddleware())
	public.POST("/verifylogin", controllers.VerifyLogin)

	protectedUsers := r.Group("/api/user", CORSMiddleware())
	protectedUsers.Use(middlewares.JWTAuthMiddleware())

	protectedUsers.OPTIONS("/porposeresidential", CORSMiddleware())
	protectedUsers.POST("/porposeresidential", controllers.ResidentialProposal)

	protectedUsers.OPTIONS("/acceptporposal", CORSMiddleware())
	protectedUsers.POST("/acceptporposal", controllers.AcceptProposal)

	protectedUsers.OPTIONS("/rejectporposal", CORSMiddleware())
	protectedUsers.POST("/rejectporposal", controllers.RejectProposal)

	protectedUsers.OPTIONS("/transfer", CORSMiddleware())
	protectedUsers.POST("/transfer", controllers.Transfer)

	protectedUsers.OPTIONS("/balance", CORSMiddleware())
	protectedUsers.GET("/balance", controllers.Balance)

	// should be admin
	protectedUsers.OPTIONS("/getallproposals", CORSMiddleware())
	protectedUsers.GET("/getallproposals", controllers.GetAllResidentialProposal)

	// should be admin
	protectedUsers.OPTIONS("/insertbill", CORSMiddleware())
	protectedUsers.POST("/insertbill", controllers.GetBill)

	protected := r.Group("/api/admin")
	protected.Use(middlewares.JWTAuthMiddleware())
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
