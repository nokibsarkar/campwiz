package database

type Evaluation struct {
	EvaluationID  string       `json:"evaluationId" gorm:"primaryKey;type:varchar(255)"`
	SubmissionID  string       `json:"submissionId" gorm:"type:varchar(255)"`
	JudgeID       string       `json:"judgeId" gorm:"type:varchar(255)"`
	ParticipantID string       `json:"participantId" gorm:"type:varchar(255)"`
	Score         int          `json:"score" gorm:"default:0"`
	Comment       string       `json:"comment" gorm:"default:null"`
	Serial        uint         `json:"serial"`
	Submission    *Submission  `json:"-"`
	Participant   *Participant `json:"-"`
	Judge         *Jury        `json:"-"`
}
