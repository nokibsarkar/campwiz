package database

type Evaluation struct {
	EvaluationID  IDType       `json:"evaluationId" gorm:"primaryKey"`
	SubmissionID  IDType       `json:"submissionId"`
	JudgeID       IDType       `json:"judgeId"`
	ParticipantID IDType       `json:"participantId"`
	Score         int          `json:"score" gorm:"default:0"`
	Comment       string       `json:"comment" gorm:"default:null"`
	Serial        uint         `json:"serial"`
	Submission    *Submission  `json:"-"`
	Participant   *Participant `json:"-"`
	Judge         *Jury        `json:"-"`
}
