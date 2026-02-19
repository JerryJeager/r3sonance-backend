package models

const (
	ArtistWeight    = 38.0
	TrackWeight     = 25.0
	RankWeight      = 15.0
	HourWeight      = 12.0
	DiversityWeight = 10.0
)

type CompatibilityResult struct {
	TotalScore        int                    `json:"total_score"`
	Percentage        float64                `json:"percentage"`
	Tier              string                 `json:"tier"`
	Breakdown         CompatibilityBreakdown `json:"breakdown"`
	SharedArtists     []SharedArtist         `json:"shared_artist"`
	SharedTracks      []SharedTrack          `json:"shared_tracks"`
	ListeningInsights ListeningInsights      `json:"listening_insights"`
}

type CompatibilityBreakdown struct {
	ArtistOverlapScore float64 `json:"artist_overlap_score"`
	TrackOverlapScore  float64 `json:"track_overlap_score"`
	RankAlignmentScore float64 `json:"rank_alignment_score"`
	ActiveHourScore    float64 `json:"active_hour_score"`
	DiversityScore     float64 `json:"diversity_score"`
}

type SharedArtist struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Image          string `json:"image"`
	RankA          int    `json:"rank_a"`
	RankB          int    `json:"rank_b"`
	RankDifference int    `json:"rank_difference"`
}

type SharedTrack struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	AlbumImage string `json:"album_image"`
	SpotifyURL string `json:"spotify_url"`
}

type ListeningInsights struct {
	MostSyncedHour string `json:"most_synced_hour"`
	ListeningTypeA string `json:"listening_type_a"`
	ListeningTypeB string `json:"listening_type_b"`
	SyncMessage    string `json:"sync_message"`
}
