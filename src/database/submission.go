package database

import "time"

type ArticleSubmission struct {
	Language   string `json:"language"`
	TotalBytes uint64 `json:"totalbytes" gorm:"default:0"`
	TotalWords uint64 `json:"totalwords" gorm:"default:0"`
	AddedBytes uint64 `json:"addedbytes" gorm:"default:0"`
	AddedWords uint64 `json:"addedwords" gorm:"default:0"`
}
type MediaSubmission struct {
	MediaType   string `json:"mediatype"`
	Height      uint64 `json:"height"`
	Width       uint64 `json:"width"`
	Duration    uint64 `json:"duration"` // in milliseconds
	Bitrate     uint64 `json:"bitrate"`  // in kbps
	Size        uint64 `json:"size"`     // in bytes
	ThumbURL    string `json:"thumburl"`
	ThumbWidth  uint64 `json:"thumbwidth"`
	ThumbHeight uint64 `json:"thumbheight"`
	License     string `json:"license"`
}
type Submission struct {
	SubmissionID  string      `json:"pageid" gorm:"primaryKey"`
	Name          string      `json:"title"`
	CampaignID    uint64      `json:"campaignId" gorm:"index"`
	URL           string      `json:"url"`
	Author        string      `json:"author"`        // The Actual Author in the Wikimedia
	AuthorID      uint64      `json:"authorId"`      // The Actual Author ID in the Wikimedia
	SubmittedByID uint64      `json:"submittedById"` // The User who submitted the article on behalf of the participant
	ParticipantID uint64      `json:"participantId" gorm:"index"`
	SubmittedAt   time.Time   `json:"submittedAt"`
	Participant   Participant `json:"-" gorm:"foreignKey:ParticipantID"`
	Submitter     User        `json:"-" gorm:"foreignKey:SubmittedByID"`
	Campaign      Campaign    `json:"-" gorm:"foreignKey:CampaignID"`
	MediaSubmission
}
