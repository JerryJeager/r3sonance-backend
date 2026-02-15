
package cmd

import (
	"log"
	"os"

	"github.com/JerryJeager/r3sonance-backend/manualwire"
	"github.com/JerryJeager/r3sonance-backend/middleware"
	"github.com/gin-gonic/gin"
)

func ExecuteApiRoutes() {
	router := gin.Default()

	router.Use(middleware.CORSMiddleware())

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Welcome",
		})
	})

	userController := manualwire.GetUserController()

	api := router.Group("/api/v1")
	users := api.Group("/users")

	users.POST("/signup", userController.CreateUser)
	users.POST("/verify-email", userController.VerifyUserEmail)
	users.POST("/login", userController.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panic("failed to run server")
	}
}

	