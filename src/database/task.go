// This would be used for running background tasks
package database

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type TaskStatus string
type TaskType string

const (
	TaskTypeImportFromCommons     TaskType = "import.commons"
	TaskTypeDistributeEvaluations TaskType = "distribute.evaluations"
)
const (
	TaskStatusPending TaskStatus = "pending"
	TaskStatusRunning TaskStatus = "running"
	TaskStatusSuccess TaskStatus = "success"
	TaskStatusFailed  TaskStatus = "failed"
)

type Task struct {
	TaskID               IDType                       `json:"taskId" gorm:"primaryKey"`
	Type                 TaskType                     `json:"type"`
	Status               TaskStatus                   `json:"status"`
	AssociatedRoundID    *IDType                      `json:"roundId" gorm:"index;null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AssociatedCampaignID *IDType                      `json:"campaignId" gorm:"index;null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	AssociatedUserID     *IDType                      `json:"userId" gorm:"index;null;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Data                 *datatypes.JSON              `json:"data"`
	CreatedAt            time.Time                    `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt            time.Time                    `json:"updatedAt" gorm:"autoUpdateTime"`
	SuccessCount         int                          `json:"successCount"`
	FailedCount          int                          `json:"failedCount"`
	FailedIds            *datatypes.JSONSlice[string] `json:"failedIds"`
	RemainingCount       int                          `json:"remainingCount"`
	CreatedByID          IDType                       `json:"createdById" gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
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
