package database

type Evaluation struct {
	ID            uint64       `json:"id" gorm:"primaryKey"`
	SubmissionID  uint64       `json:"submissionId"`
	BatchID       uint64       `json:"batchId"`
	JudgeID       uint64       `json:"judgeId"`
	ParticipantID string       `json:"participantId"`
	Score         int          `json:"score" gorm:"default:0"`
	Comment       string       `json:"comment" gorm:"default:null"`
	Serial        uint         `json:"serial"`
	Batch         *Batch       `json:"-"`
	Submission    *Submission  `json:"-" gorm:"foreignKey:SubmissionID"`
	Participant   *Participant `json:"-" gorm:"foreignKey:ParticipantID"`
	Judge         *Jury        `json:"-" gorm:"foreignKey:JudgeID"`
}
