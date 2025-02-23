package database

type Batch struct {
	BatchID uint64 `json:"batchId" gorm:"primaryKey"`
}
