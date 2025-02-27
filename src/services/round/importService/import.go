package importservice

import (
	"log"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"
	rnd "nokib/campwiz/services/round"
	"strings"

	"gorm.io/datatypes"
)

// ImportService is an interface for importing data from different sources
// All the importer services should implement this interface
type IImportSource interface {
	ImportImageResults(failedImageReason *map[string]string) ([]database.ImageResult, *map[string]string)
}
type TechnicalJudge interface {
	PreventionReason(database.ImageResult) string
}
type TaskRunner struct {
	TaskId database.IDType
	Source IImportSource
}

func NewTaskRunner(taskId database.IDType, importService IImportSource) *TaskRunner {
	return &TaskRunner{
		TaskId: taskId,
		Source: importService,
	}
}
func (b *TaskRunner) Run() (successCount, failedCount int) {
	task_repo := database.NewTaskRepository()
	round_repo := database.NewRoundRepository()
	conn, close := database.GetDB()
	defer close()

	task, err := task_repo.FindByID(conn, b.TaskId)
	if err != nil {
		log.Println("Error fetching task: ", err)
		return
	}
	defer func() {
		res := conn.Save(task)
		if res.Error != nil {
			log.Println("Error saving task: ", res.Error)
		}
	}()
	round, err := round_repo.FindByID(conn, *task.AssociatedRoundID)
	if err != nil {
		log.Println("Error fetching round: ", err)
		return
	}
	FailedImages := &map[string]string{}
	technicalJudge := rnd.NewTechnicalJudgeService(round)
	user_repo := database.NewUserRepository()
	for {
		successBatch, failedBatch := b.Source.ImportImageResults(FailedImages)

		if failedBatch != nil {
			task.FailedCount = len(*failedBatch)
			*task.FailedIds = datatypes.NewJSONType(*failedBatch)
		}
		if successBatch == nil {
			log.Println("No images found in the batch")
			break
		}
		images := []database.ImageResult{}
		for _, image := range successBatch {
			if technicalJudge.PreventionReason(image) != "" {
				images = append(images, image)
			}
		}
		task.SuccessCount += len(images)
		participants := map[string]database.IDType{}
		for _, image := range images {
			participants[image.UploaderUsername] = idgenerator.GenerateID("u")
		}
		perBatch := conn.Begin()
		username2IdMap, err := user_repo.EnsureExists(perBatch, participants)
		if err != nil {
			log.Println("Error ensuring users exist: ", err)
			perBatch.Rollback()
			task.Status = database.TaskStatusFailed
			return
		}
		submissionCount := 0
		submissions := []database.Submission{}
		for _, image := range images {
			uploaderId := username2IdMap[image.UploaderUsername]
			submission := database.Submission{
				SubmissionID:      idgenerator.GenerateID("s"),
				Name:              image.Name,
				CampaignID:        *task.AssociatedCampaignID,
				URL:               image.URL,
				Author:            image.UploaderUsername,
				SubmittedByID:     uploaderId,
				ParticipantID:     uploaderId,
				SubmittedAt:       image.SubmittedAt,
				CreatedAtExternal: &image.SubmittedAt,
				CurrentRoundID:    round.RoundID,
				MediaSubmission: database.MediaSubmission{
					MediaType:   database.MediaType(image.MediaType),
					ThumbURL:    image.URL,
					ThumbWidth:  image.Width,
					ThumbHeight: image.Height,
					License:     strings.ToUpper(image.License),
					CreditHTML:  image.CreditHTML,
					Description: image.Description,
					AudioVideoSubmission: database.AudioVideoSubmission{
						Duration: image.Duration,
						Size:     image.Size,
						Bitrate:  0,
					},
					ImageSubmission: database.ImageSubmission{
						Width:  image.Width,
						Height: image.Height,
					},
				},
			}
			submissions = append(submissions, submission)
			submissionCount++
		}
		if len(submissions) == 0 {
			// No submissions to save
			// This can happen if all the images are rejected by the technical judge
			task.Status = database.TaskStatusSuccess
			break
		}
		res := perBatch.Create(submissions)
		if res.Error != nil {
			task.Status = database.TaskStatusFailed
			log.Println("Error saving submissions: ", res.Error)
			perBatch.Rollback()
			return
		}
		*task.FailedIds = datatypes.NewJSONType(*failedBatch)
		res = perBatch.Save(task)
		if res.Error != nil {
			log.Println("Error saving task: ", res.Error)
			task.Status = database.TaskStatusFailed
			perBatch.Rollback()
			return
		}
		perBatch.Commit()
	}
	task.Status = database.TaskStatusSuccess
	return
}
