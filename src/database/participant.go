package database

type Participant struct {
	ParticipantID string `json:"id" gorm:"primaryKey"`
	Name          string `json:"name"`
}
