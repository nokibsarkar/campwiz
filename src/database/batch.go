package database

import (
	"time"

	"gorm.io/gorm"
)

type Batch struct {
	BatchID          string     `json:"batchId" gorm:"primaryKey;type:varchar(255)"`
	CreatedByID      string     `json:"createdById" gorm:"type:varchar(255)"`
	TotalSubmissions int        `json:"totalSubmissions" gorm:"default:0"`
	CreatedAt        *time.Time `json:"createdAt" gorm:"autoCreateTime"`
}
type BatchFilter struct {
	IDs   []string
	Limit int `form:"limit"`
}
type BatchRepository struct{}

func NewBatchRepository() *BatchRepository {
	return &BatchRepository{}
}
func (r *BatchRepository) CreateBatch(tx *gorm.DB, batch *Batch) (*Batch, error) {
	result := tx.Create(batch)
	return batch, result.Error
}
func (r *BatchRepository) GetBatchByID(conn *gorm.DB, batchID string) (*Batch, error) {
	batch := &Batch{}
	result := conn.First(batch, batchID)
	return batch, result.Error
}
func (r *BatchRepository) GetBatches(conn *gorm.DB, filter *BatchFilter) ([]Batch, error) {
	stmt := conn
	if filter != nil {
		if filter.Limit > 0 {
			stmt = stmt.Limit(filter.Limit)
		}
		if len(filter.IDs) > 0 {
			stmt = stmt.Where("batch_id IN ?", filter.IDs)
		}
	}
	var batches []Batch
	result := stmt.Find(&batches)
	return batches, result.Error
}
