package utils

import (
	"fmt"
	"math"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
)

func CalculateArtistScore(a, b []models.SnapshotArtist) (float64, []models.SharedArtist, float64) {

	mapA := make(map[string]int)
	mapB := make(map[string]int)

	for i, artist := range a {
		mapA[artist.ID] = i + 1
	}

	for i, artist := range b {
		mapB[artist.ID] = i + 1
	}

	intersection := 0
	union := make(map[string]bool)
	var shared []models.SharedArtist
	var rankSimilarities []float64

	for id := range mapA {
		union[id] = true
	}
	for id := range mapB {
		union[id] = true
	}

	for id, rankA := range mapA {
		if rankB, ok := mapB[id]; ok {
			intersection++

			diff := math.Abs(float64(rankA - rankB))
			rankSim := 1 - (diff / 50.0)
			rankSimilarities = append(rankSimilarities, rankSim)

			shared = append(shared, models.SharedArtist{
				ID:             id,
				Name:           findArtistName(a, id),
				Image:          findArtistImage(a, id),
				RankA:          rankA,
				RankB:          rankB,
				RankDifference: int(diff),
			})
		}
	}

	if len(union) == 0 {
		return 0, nil, 0
	}

	jaccard := float64(intersection) / float64(len(union))
	artistScore := jaccard * models.ArtistWeight

	avgRankSim := 0.0
	if len(rankSimilarities) > 0 {
		sum := 0.0
		for _, v := range rankSimilarities {
			sum += v
		}
		avgRankSim = sum / float64(len(rankSimilarities))
	}

	return artistScore, shared, avgRankSim
}

func CalculateTrackScore(a, b []models.SnapshotTrack) (float64, []models.SharedTrack) {

	setA := make(map[string]models.SnapshotTrack)
	setB := make(map[string]models.SnapshotTrack)

	for _, t := range a {
		setA[t.ID] = t
	}
	for _, t := range b {
		setB[t.ID] = t
	}

	union := make(map[string]bool)
	intersection := 0
	var shared []models.SharedTrack

	for id := range setA {
		union[id] = true
	}
	for id := range setB {
		union[id] = true
	}

	for id, trackA := range setA {
		if trackB, ok := setB[id]; ok {
			intersection++

			shared = append(shared, models.SharedTrack{
				ID:         id,
				Name:       trackA.Name,
				AlbumImage: trackA.AlbumImage,
				SpotifyURL: trackB.SpotifyURL,
			})
		}
	}

	if len(union) == 0 {
		return 0, nil
	}

	jaccard := float64(intersection) / float64(len(union))
	return jaccard * models.TrackWeight, shared
}

func CalculateActiveHourScore(a, b map[int]int) (float64, models.ListeningInsights) {

	insights := models.ListeningInsights{}

	if len(a) < 2 || len(b) < 2 {
		insights.SyncMessage = "Not enough listening data to analyze patterns."
		return models.HourWeight / 2, insights
	}

	vecA := make([]float64, 24)
	vecB := make([]float64, 24)

	totalA := 0
	totalB := 0

	for _, v := range a {
		totalA += v
	}
	for _, v := range b {
		totalB += v
	}

	for hour, count := range a {
		vecA[hour] = float64(count) / float64(totalA)
	}
	for hour, count := range b {
		vecB[hour] = float64(count) / float64(totalB)
	}

	dot := 0.0
	magA := 0.0
	magB := 0.0

	for i := 0; i < 24; i++ {
		dot += vecA[i] * vecB[i]
		magA += vecA[i] * vecA[i]
		magB += vecB[i] * vecB[i]
	}

	if magA == 0 || magB == 0 {
		insights.SyncMessage = "Listening data too sparse to compare."
		return 0, insights
	}

	cosine := dot / (math.Sqrt(magA) * math.Sqrt(magB))
	score := cosine * models.HourWeight

	peakA := getPeakHour(a)
	peakB := getPeakHour(b)

	insights.ListeningTypeA = categorizeHour(peakA)
	insights.ListeningTypeB = categorizeHour(peakB)

	if peakA == peakB {
		insights.MostSyncedHour = formatHour(peakA)
		insights.SyncMessage = "You both peak around " + insights.MostSyncedHour
	} else {
		insights.SyncMessage = "You have different listening rhythms."
	}

	return score, insights
}

func CalculateDiversityScore(a, b int) float64 {

	if a == 0 && b == 0 {
		return models.DiversityWeight
	}

	max := float64(math.Max(float64(a), float64(b)))
	diff := math.Abs(float64(a - b))

	similarity := 1 - (diff / max)
	return similarity * models.DiversityWeight
}

func GetTier(score float64) string {
	switch {
	case score >= 90:
		return "Soul Synced"
	case score >= 75:
		return "Highly in Tune"
	case score >= 60:
		return "Strong Vibe"
	case score >= 40:
		return "Some Overlap"
	default:
		return "Different Wavelengths"
	}
}

func findArtistName(artists []models.SnapshotArtist, id string) string {
	for _, a := range artists {
		if a.ID == id {
			return a.Name
		}
	}
	return ""
}

func findArtistImage(artists []models.SnapshotArtist, id string) string {
	for _, a := range artists {
		if a.ID == id {
			return a.Image
		}
	}
	return ""
}

func getPeakHour(hours map[int]int) int {
	max := -1
	peak := 0
	for h, v := range hours {
		if v > max {
			max = v
			peak = h
		}
	}
	return peak
}

func categorizeHour(hour int) string {
	switch {
	case hour >= 5 && hour <= 11:
		return "Morning Listener"
	case hour >= 12 && hour <= 17:
		return "Afternoon Listener"
	case hour >= 18 && hour <= 21:
		return "Evening Listener"
	default:
		return "Night Owl"
	}
}

func formatHour(hour int) string {
	return fmt.Sprintf("%02d:00", hour)
}
