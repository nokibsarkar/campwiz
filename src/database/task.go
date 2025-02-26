// This would be used for running background tasks
package database

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Task struct {
	TaskID               IDType                       `json:"taskId" gorm:"primaryKey"`
	Type                 string                       `json:"type"`
	Status               string                       `json:"status"`
	AssociatedRoundID    *IDType                      `json:"roundId" gorm:"index;null"`
	AssociatedCampaignID *IDType                      `json:"campaignId" gorm:"index;null"`
	AssociatedUserID     *IDType                      `json:"userId" gorm:"index;null"`
	Data                 *datatypes.JSON              `json:"data" gorm:"type:json"`
	CreatedAt            string                       `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            string                       `json:"updatedAt" gorm:"autoUpdateTime"`
	SuccessCount         int                          `json:"successCount"`
	FailedCount          int                          `json:"failedCount"`
	FailedIds            *datatypes.JSONSlice[string] `json:"failedIds" gorm:"type:json"`
	RemainingCount       int                          `json:"remainingCount"`
	CreatedByID          IDType                       `json:"createdById"`
	Submittor            User                         `json:"-" gorm:"foreignKey:CreatedByID;references:UserID"`
}

type TaskRepository struct{}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{}
}

func (r *TaskRepository) Create(tx *gorm.DB, task *Task) (*Task, error) {
	err := tx.Create(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (r *TaskRepository) FindByID(tx *gorm.DB, taskId IDType) (*Task, error) {
	task := &Task{TaskID: taskId}
	err := tx.Where(task).First(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}
func (r *TaskRepository) Update(tx *gorm.DB, task *Task) (*Task, error) {
	err := tx.Save(task).Error
	if err != nil {
		return nil, err
	}
	return task, nil
}
