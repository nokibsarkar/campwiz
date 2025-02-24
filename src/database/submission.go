package database

import (
	"time"

	"gorm.io/gorm"
)

type ArticleSubmission struct {
	Language   string `json:"language"`
	TotalBytes uint64 `json:"totalbytes" gorm:"default:0"`
	TotalWords uint64 `json:"totalwords" gorm:"default:0"`
	AddedBytes uint64 `json:"addedbytes" gorm:"default:0"`
	AddedWords uint64 `json:"addedwords" gorm:"default:0"`
}
type ImageSubmission struct {
	Width  uint64 `json:"width"`
	Height uint64 `json:"height"`
}
type AudioVideoSubmission struct {
	Duration uint64 `json:"duration"` // in milliseconds
	Bitrate  uint64 `json:"bitrate"`  // in kbps
	Size     uint64 `json:"size"`     // in bytes
}
type MediaSubmission struct {
	MediaType   string  `json:"mediatype"`
	ThumbURL    string  `json:"thumburl"`
	ThumbWidth  uint64  `json:"thumbwidth"`
	ThumbHeight uint64  `json:"thumbheight"`
	License     string  `json:"license"`
	BatchID     *string `json:"batchId" gorm:"null;varchar(255)"`
	Batch       *Batch  `json:"-" gorm:"foreignKey:BatchID"`
	ImageSubmission
	AudioVideoSubmission
}
type Submission struct {
	SubmissionID  string      `json:"pageid" `
	Name          string      `json:"title"`
	CampaignID    uint64      `json:"campaignId" gorm:"index"`
	URL           string      `json:"url"`
	Author        string      `json:"author"`        // The Actual Author in the Wikimedia
	AuthorID      uint64      `json:"authorId"`      // The Actual Author ID in the Wikimedia
	SubmittedByID uint64      `json:"submittedById"` // The User who submitted the article on behalf of the participant
	ParticipantID uint64      `json:"participantId"`
	SubmittedAt   time.Time   `json:"submittedAt"`
	Participant   Participant `json:"-"`
	Submitter     User        `json:"-" gorm:"foreignKey:SubmittedByID"`
	Campaign      *Campaign   `json:"-" gorm:"foreignKey:CampaignID"`
	MediaSubmission
}
type SubmissionRepository struct{}

func NewSubmissionRepository() *SubmissionRepository {
	return &SubmissionRepository{}
}
func (r *SubmissionRepository) CreateSubmission(tx *gorm.DB, submission *Submission) error {
	result := tx.Create(submission)
	return result.Error
}
