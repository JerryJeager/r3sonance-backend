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

	users.GET("/spotify/login", userController.SpotifyLogin)
	users.GET("/spotify/callback", userController.SpotifyCallback)

	users.Use(middleware.SpotfiyAuthMiddleware())
	{
		users.GET("/spotify/snapshot", userController.GetUserMusicSnapshot)
		users.GET("/profile", userController.GetUser)
		users.GET("/compatibility/:public_id", userController.GetMusicCompatibility)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Panic("failed to run server")
	}
}
