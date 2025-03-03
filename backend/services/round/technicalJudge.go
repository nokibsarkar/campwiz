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
		MinimumUploadDate: round.StartDate,
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
func (j *TechnicalJudgeService) PreventionReason(img database.ImageResult) string {
	if !j.MinimumUploadDate.IsZero() && img.SubmittedAt.Before(j.MinimumUploadDate) {
		// log.Printf("Image %s is not allowed because it was uploaded before %s", img.Name, j.MinimumUploadDate)
		return "before-minimum-upload-date"
	}
	if img.Resolution < j.MinimumResolution {
		// log.Printf("Image %s is not allowed because it has a resolution of %d which is less than %d", img.Name, img.Resolution, j.MinimumResolution)
		return "below-minimum-resolution"
	}
	if img.Size < j.MinimumSize {
		// log.Printf("Image %s is not allowed because it has a size of %d which is less than %d", img.Name, img.Size, j.MinimumSize)
		return "below-minimum-size"
	}
	if j.AllowedTypes != nil && !j.AllowedTypes.Contains(database.MediaType(img.MediaType)) {
		// log.Printf("Image %s is not allowed because it is of type %s", img.Name, img.MediaType)
		return "not-allowed-type"
	}
	return ""
}
