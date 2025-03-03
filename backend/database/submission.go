package database

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SubmissionListFilter struct {
	CampaignID    IDType `form:"campaignId"`
	RoundID       IDType `form:"roundId"`
	ParticipantID IDType `form:"participantId"`
	CommonFilter
}
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
	MediaType   MediaType      `json:"mediatype" gorm:"not null;default:'BITMAP'"`
	ThumbURL    string         `json:"thumburl"`
	ThumbWidth  uint64         `json:"thumbwidth"`
	ThumbHeight uint64         `json:"thumbheight"`
	License     string         `json:"license"`
	Description string         `json:"description"`
	CreditHTML  string         `json:"creditHTML"`
	Metadata    datatypes.JSON `json:"metadata" gorm:"type:json"`
	ImageSubmission
	AudioVideoSubmission
}
type Submission struct {
	SubmissionID IDType `json:"pageid" gorm:"primaryKey"`
	Name         string `json:"title"`
	CampaignID   IDType `json:"campaignId" gorm:"null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	URL          string `json:"url"`
	// The Actual Author in the Wikimedia
	Author UserName `json:"author"`
	// The User who submitted the article on behalf of the participant
	SubmittedByID  IDType    `json:"submittedById" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParticipantID  IDType    `json:"participantId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CurrentRoundID IDType    `json:"currentRoundId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	SubmittedAt    time.Time `json:"submittedAt" gorm:"type:datetime"`
	Participant    User      `json:"-" gorm:"foreignKey:ParticipantID"`
	Submitter      User      `json:"-" gorm:"foreignKey:SubmittedByID"`
	// Campaign          *Campaign  `json:"-" gorm:"foreignKey:CampaignID"`
	CreatedAtExternal *time.Time `json:"createdAtServer"`
	CurrentRound      *Round     `json:"-" gorm:"foreignKey:CurrentRoundID"`
	MediaSubmission
}
type SubmissionSelectID struct {
	SubmissionID IDType
}
type SubmissionRepository struct{}

func NewSubmissionRepository() *SubmissionRepository {
	return &SubmissionRepository{}
}
func (r *SubmissionRepository) CreateSubmission(tx *gorm.DB, submission *Submission) error {
	result := tx.Create(submission)
	return result.Error
}
func (r *SubmissionRepository) FindSubmissionByID(tx *gorm.DB, submissionID IDType) (*Submission, error) {
	submission := &Submission{}
	result := tx.First(submission, &Submission{SubmissionID: submissionID})
	return submission, result.Error
}
func (r *SubmissionRepository) ListAllSubmissions(tx *gorm.DB, filter *SubmissionListFilter) ([]Submission, error) {
	var submissions []Submission
	condition := &Submission{}
	if filter != nil {
		if filter.CampaignID != "" {
			condition.CampaignID = filter.CampaignID
		}
		if filter.RoundID != "" {
			condition.CurrentRoundID = filter.RoundID
		}
		if filter.ParticipantID != "" {
			condition.ParticipantID = filter.ParticipantID
		}
	}
	where := tx.Where(condition)
	if filter.ContinueToken != "" {
		where = where.Where("submission_id > ?", filter.ContinueToken)
	}

	stmt := where //.Order("submission_id")
	if filter.Limit > 0 {
		stmt = stmt.Limit(max(100, filter.Limit))
	}
	result := stmt.Find(&submissions)
	return submissions, result.Error
}
func (r *SubmissionRepository) GetSubmissionCount(tx *gorm.DB, filter *SubmissionListFilter) (int64, error) {
	condition := &Submission{}
	if filter != nil {
		if filter.CampaignID != "" {
			condition.CampaignID = filter.CampaignID
		}
		if filter.RoundID != "" {
			condition.CurrentRoundID = filter.RoundID
		}
		if filter.ParticipantID != "" {
			condition.ParticipantID = filter.ParticipantID
		}
	}
	var count int64
	result := tx.Model(&Submission{}).Where(condition).Count(&count)
	return count, result.Error
}
func (r *SubmissionRepository) ListAllSubmissionIDs(tx *gorm.DB, filter *SubmissionListFilter) ([]SubmissionSelectID, error) {
	var submissionIDs []SubmissionSelectID
	condition := &Submission{}
	if filter != nil {
		if filter.CampaignID != "" {
			condition.CampaignID = filter.CampaignID
		}
		if filter.RoundID != "" {
			condition.CurrentRoundID = filter.RoundID
		}
		if filter.ParticipantID != "" {
			condition.ParticipantID = filter.ParticipantID
		}
	}
	where := tx.Where(condition)
	if filter.ContinueToken != "" {
		where = where.Where("submission_id > ?", filter.ContinueToken)
	}

	stmt := where //.Order("submission_id")
	if filter.Limit > 0 {
		stmt = stmt.Limit(max(100, filter.Limit))
	}
	result := stmt.Model(&Submission{}).Find(&submissionIDs)
	return submissionIDs, result.Error
}
