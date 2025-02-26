package round

import (
	"nokib/campwiz/database"
	"time"
)

type TechnicalJudgeService struct {
	AllowedTypes      database.MediaTypeSet
	MinimumUploadDate time.Time
	MinimumResolution uint64
	MinimumSize       uint64
	// This would be a list of persons who are not allowed to submit images
	// Thes include the banned users, judges, coordinators, moderators etc
	Blacklist []string
}

func NewTechnicalJudgeService(round *database.Round) *TechnicalJudgeService {
	return &TechnicalJudgeService{
		AllowedTypes:      round.AllowedMediaTypes,
		MinimumUploadDate: time.Now().AddDate(0, 0, -1),
		MinimumResolution: uint64(round.MinimumResolution),
		MinimumSize:       uint64(round.MinimumTotalBytes),
		Blacklist:         []string{},
	}
}

// This method would perform some basic checks to see if the image is prevented from submission
// It would consider
// - For images
//   - Minimum Upload Date
//   - Minimum Resolution
//   - Minimum Size
//   - Whether Image allowed or not
func (j *TechnicalJudgeService) IsMediaAllowed(img database.ImageResult) bool {
	return true
	if img.SubmittedAt.Before(j.MinimumUploadDate) {
		return false
	}
	// if img. < j.MinimumResolution {
	// 	return false
	// }
	if img.Size < j.MinimumSize {
		return false
	}
	return j.AllowedTypes.Contains(database.MediaType(img.MediaType))
}
