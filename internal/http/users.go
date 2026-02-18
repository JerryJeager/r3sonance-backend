package http

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"github.com/JerryJeager/r3sonance-backend/internal/service/spotify"
	"github.com/JerryJeager/r3sonance-backend/internal/service/users"
	"github.com/JerryJeager/r3sonance-backend/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type UserController struct {
	serv users.UserSv
}

func NewUserController(serv users.UserSv) *UserController {
	return &UserController{serv: serv}
}

func (c *UserController) SpotifyLogin(ctx *gin.Context) {

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	state, err := utils.GenerateState()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate state"})
		return
	}

	ctx.SetCookie("spotify_auth_state", state, 600, "/", "127.0.0.1", false, true)

	baseURL := "https://accounts.spotify.com/authorize"

	params := url.Values{}
	params.Add("client_id", clientID)
	params.Add("response_type", "code")
	params.Add("redirect_uri", redirectURI)
	params.Add("scope", "user-top-read user-read-recently-played user-read-email user-read-private playlist-read-private")
	params.Add("state", state)

	authURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	ctx.Redirect(http.StatusFound, authURL)
}

func (c *UserController) SpotifyCallback(ctx *gin.Context) {

	code := ctx.Query("code")
	state := ctx.Query("state")

	if code == "" || state == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing code or state"})
		return
	}

	// Validate state
	storedState, err := ctx.Cookie("spotify_auth_state")
	if err != nil || storedState != state {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	redirectURI := os.Getenv("SPOTIFY_REDIRECT_URI")

	restyClient := resty.New()

	// Exchange code for token
	var tokenResp models.SpotifyTokenResponse

	resp, err := restyClient.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBasicAuth(clientID, clientSecret).
		SetFormData(map[string]string{
			"grant_type":   "authorization_code",
			"code":         code,
			"redirect_uri": redirectURI,
		}).
		SetResult(&tokenResp).
		Post("https://accounts.spotify.com/api/token")

	if err != nil || resp.StatusCode() != 200 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}

	var profile models.SpotifyProfile

	resp, err = restyClient.R().
		SetHeader("Authorization", "Bearer "+tokenResp.AccessToken).
		SetResult(&profile).
		Get("https://api.spotify.com/v1/me")

	if err != nil || resp.StatusCode() != 200 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch profile"})
		return
	}

	err = c.serv.CreateUser(ctx, &profile, &tokenResp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	go func() {
		// if shouldUpdate, _ := c.serv.ShouldUpdateUserMusicSnapshot(ctx, profile.Email); shouldUpdate {
		spotifyClient := spotify.NewSpotifyClient(tokenResp.AccessToken)
		snapshot := spotify.GetUserMusicSnapshot(spotifyClient)
		if err := c.serv.CreateUserMusicSnapshot(ctx, profile.Email, snapshot); err != nil {
			log.Printf("failed to update user music snapshot: %s\n", err.Error())
		}
		// }
	}()

	ctx.JSON(http.StatusOK, gin.H{
		"message": "spotify login successful",
		"user":    profile,
	})
}

func (c *UserController) GetUser(ctx *gin.Context) {
	email, exists := ctx.Get("user_email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found"})
		return
	}

	user, err := c.serv.GetUserByEmail(ctx, email.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"email":     user.Email,
		"name":      user.DisplayName,
		"country":   user.Country,
		"public_id": user.PublicID,
	})
}

func (c *UserController) GetUserMusicSnapshot(ctx *gin.Context) {
	email, exists := ctx.Get("user_email")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user email not found"})
		return
	}

	snapshot, err := c.serv.GetUserMusicSnapshot(ctx, email.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user music snapshot"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"snapshot": gin.H{
			"top_artists": snapshot.TopArtists,
			"top_tracks":  snapshot.TopTracks,
		},
	})
}
