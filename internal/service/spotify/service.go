package spotify

import (
	"fmt"
	"strings"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"github.com/JerryJeager/r3sonance-backend/internal/utils"
	"github.com/go-resty/resty/v2"
)

type SpotifyClient struct {
	client      *resty.Client
	accessToken string
}

func NewSpotifyClient(token string) *SpotifyClient {
	c := resty.New()

	return &SpotifyClient{
		client:      c,
		accessToken: token,
	}
}

func (s *SpotifyClient) Get(endpoint string, query map[string]string, result interface{}) error {
	resp, err := s.client.R().
		SetHeader("Authorization", "Bearer "+s.accessToken).
		SetQueryParams(query).
		SetResult(result).
		Get("https://api.spotify.com/v1" + endpoint)

	if err != nil || resp.StatusCode() != 200 {
		return fmt.Errorf("spotify api error")
	}
	return nil
}

func (s *SpotifyClient) GetTopArtists() (*models.TopArtistsResponse, error) {
	var result models.TopArtistsResponse
	err := s.Get("/me/top/artists", map[string]string{
		"limit":      "50",
		"time_range": "medium_term",
	}, &result)
	return &result, err
}

func (s *SpotifyClient) GetTopTracks() (*models.TopTracksResponse, error) {
	var result models.TopTracksResponse

	err := s.Get("/me/top/tracks", map[string]string{
		"limit":      "50",
		"time_range": "medium_term",
	}, &result)

	return &result, err
}

func (s *SpotifyClient) GetAudioFeatures(ids []string) (*models.AudioFeaturesResponse, error) {
	var result models.AudioFeaturesResponse

	err := s.Get("/audio-features", map[string]string{
		"ids": strings.Join(ids, ","),
	}, &result)

	return &result, err
}

func (s *SpotifyClient) GetRecentlyPlayed() (*models.RecentlyPlayedResponse, error) {
	var result models.RecentlyPlayedResponse

	err := s.Get("/me/player/recently-played", map[string]string{
		"limit": "50",
	}, &result)

	return &result, err
}

func GetUserMusicSnapshot(spotifyClient *SpotifyClient) *models.BuiltUserMusicSnapchat {
	topArtistsResp, _ := spotifyClient.GetTopArtists()
	topTracksResp, _ := spotifyClient.GetTopTracks()

	// trackIDs := utils.ExtractTrackIDs(topTracksResp)
	// audioResp, _ := spotifyClient.GetAudioFeatures(trackIDs)

	recentResp, _ := spotifyClient.GetRecentlyPlayed()

	topArtists := utils.TransformTopArtists(topArtistsResp)
	topTracks := utils.TransformTopTracks(topTracksResp)
	// audioProfile := utils.ComputeAudioProfile(audioResp)
	listeningPattern := utils.ComputeListeningPattern(recentResp)
	// topGenres := utils.ExtractTopGenres(topArtistsResp)

	return &models.BuiltUserMusicSnapchat{
		TopArtists:       topArtists,
		TopTracks:        topTracks,
		ListeningPattern: listeningPattern,
	}
}
