package compatibility

import (
	"math"

	"github.com/JerryJeager/r3sonance-backend/internal/models"
	"github.com/JerryJeager/r3sonance-backend/internal/utils"
)

func CalculateCompatibility(a, b models.UserMusicSnapshot) models.CompatibilityResult {

	artistScore, sharedArtists, avgRankSim :=
		utils.CalculateArtistScore(a.TopArtists, b.TopArtists)

	trackScore, sharedTracks :=
		utils.CalculateTrackScore(a.TopTracks, b.TopTracks)

	rankScore := avgRankSim * models.RankWeight

	hourScore, insights :=
		utils.CalculateActiveHourScore(a.ListeningPattern.ActiveHours,
			b.ListeningPattern.ActiveHours)

	diversityScore :=
		utils.CalculateDiversityScore(
			a.ListeningPattern.UniqueArtists,
			b.ListeningPattern.UniqueArtists)

	total := artistScore +
		trackScore +
		rankScore +
		hourScore +
		diversityScore

	if total > 100 {
		total = 100
	}

	return models.CompatibilityResult{
		TotalScore: int(math.Round(total)),
		Percentage: math.Round(total),
		Tier:       utils.GetTier(total),

		Breakdown: models.CompatibilityBreakdown{
			ArtistOverlapScore: artistScore,
			TrackOverlapScore:  trackScore,
			RankAlignmentScore: rankScore,
			ActiveHourScore:    hourScore,
			DiversityScore:     diversityScore,
		},

		SharedArtists:     sharedArtists,
		SharedTracks:      sharedTracks,
		ListeningInsights: insights,
	}
}
