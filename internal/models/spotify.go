package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type UserMusicSnapshot struct {
	ID         uuid.UUID          `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	TopArtists SnapshotArtistList `json:"top_artists"`
	TopTracks  SnapshotTrackList  `json:"top_tracks"`
	// TopGenres        SnapshotGenreList        `json:"top_genres"`
	// AudioProfile     SnapshotAudioProfile     `json:"audio_profile"`
	ListeningPattern SnapshotListeningPattern `json:"listening_pattern"`
	CreatedAt        time.Time                `json:"created_at"`
	UpdatedAt        time.Time                `json:"updated_at"`
}

type SnapshotArtist struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	Genres     []string `json:"genres"`
	Popularity int      `json:"popularity"`
	Image      string   `json:"image"`
	SpotifyURL string   `json:"spotify_url"`
}

type SnapshotArtistRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SnapshotTrack struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Artists    []SnapshotArtistRef `json:"artists"`
	AlbumName  string              `json:"album_name"`
	AlbumImage string              `json:"album_image"`
	Popularity int                 `json:"popularity"`
	DurationMs int                 `json:"duration_ms"`
	SpotifyURL string              `json:"spotify_url"`
}

type SnapshotAudioProfile struct {
	Averages SnapshotAudioAverages `json:"averages"`
}

type SnapshotAudioAverages struct {
	Danceability     float64 `json:"danceability"`
	Energy           float64 `json:"energy"`
	Valence          float64 `json:"valence"`
	Tempo            float64 `json:"tempo"`
	Acousticness     float64 `json:"acousticness"`
	Instrumentalness float64 `json:"instrumentalness"`
	Liveness         float64 `json:"liveness"`
	Speechiness      float64 `json:"speechiness"`
}

type SnapshotListeningPattern struct {
	UniqueArtists      int         `json:"unique_artists"`
	RecentCount        int         `json:"recent_count"`
	AvgTrackDurationMs int         `json:"avg_track_duration_ms"`
	ActiveHours        map[int]int `json:"active_hours"`
	WeekendRatio       float64     `json:"weekend_ratio"`
	RepeatRatio        float64     `json:"repeat_ratio"`
}

type SnapshotArtistList []SnapshotArtist
type SnapshotTrackList []SnapshotTrack
type SnapshotGenreList []string

func (a SnapshotArtistList) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *SnapshotArtistList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SnapshotArtistList")
	}
	return json.Unmarshal(bytes, a)
}

func (p SnapshotAudioProfile) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *SnapshotAudioProfile) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SnapshotAudioProfile")
	}
	return json.Unmarshal(bytes, p)
}

func (p *SnapshotGenreList) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *SnapshotGenreList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SnapshotGenreList")
	}
	return json.Unmarshal(bytes, p)
}

func (p *SnapshotListeningPattern) Value() (driver.Value, error) {
	return json.Marshal(p)
}

func (p *SnapshotListeningPattern) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SnapshotListeningPattern")
	}
	return json.Unmarshal(bytes, p)
}

func (t SnapshotTrackList) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *SnapshotTrackList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan SnapshotTrackList")
	}
	return json.Unmarshal(bytes, t)
}

type TopArtistsResponse struct {
	Items []struct {
		ID         string   `json:"id"`
		Name       string   `json:"name"`
		Genres     []string `json:"genres"`
		Popularity int      `json:"popularity"`
		Images     []struct {
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"images"`
		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"items"`
}

type TopTracksResponse struct {
	Items []struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
		DurationMs int    `json:"duration_ms"`

		Artists []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"artists"`

		Album struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			Images []struct {
				URL string `json:"url"`
			} `json:"images"`
		} `json:"album"`

		ExternalURLs struct {
			Spotify string `json:"spotify"`
		} `json:"external_urls"`
	} `json:"items"`
}

type AudioFeaturesResponse struct {
	AudioFeatures []struct {
		ID               string  `json:"id"`
		Danceability     float64 `json:"danceability"`
		Energy           float64 `json:"energy"`
		Valence          float64 `json:"valence"`
		Tempo            float64 `json:"tempo"`
		Acousticness     float64 `json:"acousticness"`
		Instrumentalness float64 `json:"instrumentalness"`
		Liveness         float64 `json:"liveness"`
		Speechiness      float64 `json:"speechiness"`
	} `json:"audio_features"`
}

type RecentlyPlayedResponse struct {
	Items []struct {
		PlayedAt string `json:"played_at"`

		Track struct {
			ID         string `json:"id"`
			Name       string `json:"name"`
			DurationMs int    `json:"duration_ms"`
			Popularity int    `json:"popularity"`

			Artists []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"artists"`

			Album struct {
				ID     string `json:"id"`
				Name   string `json:"name"`
				Images []struct {
					URL    string `json:"url"`
					Height int    `json:"height"`
					Width  int    `json:"width"`
				} `json:"images"`
			} `json:"album"`

			ExternalURLs struct {
				Spotify string `json:"spotify"`
			} `json:"external_urls"`
		} `json:"track"`
	} `json:"items"`

	Next    string `json:"next"`
	Cursors struct {
		After  string `json:"after"`
		Before string `json:"before"`
	} `json:"cursors"`

	Limit int `json:"limit"`
}

type BuiltUserMusicSnapchat struct {
	TopArtists       SnapshotArtistList
	TopTracks        SnapshotTrackList
	AudioProfile     SnapshotAudioProfile
	ListeningPattern SnapshotListeningPattern
	TopGenre         []string
}
