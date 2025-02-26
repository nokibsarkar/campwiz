// This would be used for running background tasks
package database

import (
	"time"

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
	Data                 *datatypes.JSON              `json:"data"`
	CreatedAt            time.Time                    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            time.Time                    `json:"updatedAt" gorm:"autoUpdateTime"`
	SuccessCount         int                          `json:"successCount"`
	FailedCount          int                          `json:"failedCount"`
	FailedIds            *datatypes.JSONSlice[string] `json:"failedIds"`
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
	task := &Task{}
	err := tx.Find(task, &Task{TaskID: taskId}).First(task).Error
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
