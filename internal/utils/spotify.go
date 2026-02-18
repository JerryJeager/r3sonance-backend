package utils

import (
	"sort"
	"time"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
)

func TransformTopTracks(resp *models.TopTracksResponse) models.SnapshotTrackList {
	var tracks models.SnapshotTrackList

	for _, t := range resp.Items {

		albumImage := ""
		if len(t.Album.Images) > 0 {
			albumImage = t.Album.Images[0].URL
		}

		var artists []models.SnapshotArtistRef
		for _, a := range t.Artists {
			artists = append(artists, models.SnapshotArtistRef{
				ID:   a.ID,
				Name: a.Name,
			})
		}

		tracks = append(tracks, models.SnapshotTrack{
			ID:         t.ID,
			Name:       t.Name,
			Artists:    artists,
			AlbumName:  t.Album.Name,
			AlbumImage: albumImage,
			Popularity: t.Popularity,
			DurationMs: t.DurationMs,
			SpotifyURL: t.ExternalURLs.Spotify,
		})
	}

	return tracks
}

func TransformTopArtists(resp *models.TopArtistsResponse) models.SnapshotArtistList {
	var artists models.SnapshotArtistList

	for _, a := range resp.Items {
		image := ""
		if len(a.Images) > 0 {
			image = a.Images[0].URL
		}

		artists = append(artists, models.SnapshotArtist{
			ID:         a.ID,
			Name:       a.Name,
			Genres:     a.Genres,
			Popularity: a.Popularity,
			Image:      image,
			SpotifyURL: a.ExternalURLs.Spotify,
		})
	}

	return artists
}

func ExtractTrackIDs(resp *models.TopTracksResponse) []string {
	var ids []string
	for _, t := range resp.Items {
		ids = append(ids, t.ID)
	}
	return ids
}

func ComputeAudioProfile(resp *models.AudioFeaturesResponse) models.SnapshotAudioProfile {

	var totalDance, totalEnergy, totalValence, totalTempo float64
	count := 0

	for _, f := range resp.AudioFeatures {
		if f.ID == "" {
			continue
		}

		totalDance += f.Danceability
		totalEnergy += f.Energy
		totalValence += f.Valence
		totalTempo += f.Tempo
		count++
	}

	if count == 0 {
		return models.SnapshotAudioProfile{}
	}

	return models.SnapshotAudioProfile{
		Averages: models.SnapshotAudioAverages{
			Danceability: totalDance / float64(count),
			Energy:       totalEnergy / float64(count),
			Valence:      totalValence / float64(count),
			Tempo:        totalTempo / float64(count),
		},
	}
}

func ComputeListeningPattern(resp *models.RecentlyPlayedResponse) models.SnapshotListeningPattern {

	hourCount := make(map[int]int)
	artistSet := make(map[string]bool)

	for _, item := range resp.Items {
		t, _ := time.Parse(time.RFC3339, item.PlayedAt)
		hourCount[t.Hour()]++
		artistSet[item.Track.Artists[0].ID] = true
	}

	return models.SnapshotListeningPattern{
		UniqueArtists: len(artistSet),
		ActiveHours:   hourCount,
	}
}

func ExtractTopGenres(artists *models.TopArtistsResponse) []string {
	genreCount := make(map[string]int)

	for _, a := range artists.Items {
		for _, g := range a.Genres {
			genreCount[g]++
		}
	}

	type pair struct {
		Genre string
		Count int
	}

	var sorted []pair
	for g, c := range genreCount {
		sorted = append(sorted, pair{g, c})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Count > sorted[j].Count
	})

	var topGenres []string
	for i := 0; i < len(sorted) && i < 10; i++ {
		topGenres = append(topGenres, sorted[i].Genre)
	}

	return topGenres
}
