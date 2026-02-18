package middleware

import (
	"fmt"
	"net/http"

	auth "github.com/JerryJeager/r3sonance-backend/internal/http"
	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := auth.ValidateToken(c)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     "Bad request",
				"message":    "Authentication failed",
				"statusCode": http.StatusUnauthorized,
			})
			fmt.Println(err)
			c.Abort()
			return
		}
		c.Set("user_id", id)

		c.Next()
	}
}

func SpotfiyAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		restyClient := resty.New()
		token := auth.GetTokenFromRequest(c)
		var profile models.SpotifyProfile

		resp, err := restyClient.R().
			SetHeader("Authorization", "Bearer "+token).
			SetResult(&profile).
			Get("https://api.spotify.com/v1/me")

		if err != nil || resp.StatusCode() != 200 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":     "Bad request",
				"message":    "Authentication failed",
				"statusCode": http.StatusUnauthorized,
			})
			fmt.Println(err)
			c.Abort()
			return
		}

		c.Set("user_email", profile.Email)
		c.Next()
	}
}
