package database

type Batch struct {
	ID uint64 `json:"id" gorm:"primaryKey"`
}
