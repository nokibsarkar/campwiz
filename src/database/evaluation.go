package database

type Evaluation struct {
	EvaluationID  IDType      `json:"evaluationId" gorm:"primaryKey"`
	SubmissionID  IDType      `json:"submissionId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	JudgeID       IDType      `json:"judgeId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ParticipantID IDType      `json:"participantId" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Score         int         `json:"score" gorm:"default:0"`
	Comment       string      `json:"comment" gorm:"default:null"`
	Serial        uint        `json:"serial"`
	Submission    *Submission `json:"-"`
	Participant   *User       `json:"-" gorm:"foreignKey:ParticipantID"`
	Judge         *Jury       `json:"-"`
}
