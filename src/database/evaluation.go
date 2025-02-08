package database

type Evaluation struct {
	ID           uint64 `json:"id" gorm:"primaryKey"`
	SubmissionID uint64 `json:"submissionId"`
	BatchID      uint64 `json:"batchId"`
	JudgeID      uint64 `json:"judgeId"`
	Score        int    `json:"score"`
	Comment      string `json:"comment" gorm:"default:null"`
	Serial       uint   `json:"serial"`
	Batch        *Batch `json:"-"`
}
