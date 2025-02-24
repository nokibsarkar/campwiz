package database

type Participant struct {
	ParticipantID string `json:"id" gorm:"primaryKey;type:varchar(255)"`
	Name          string `json:"name"`
}
