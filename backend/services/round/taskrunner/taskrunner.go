package taskrunner

import (
	"log"
	"nokib/campwiz/database"
	idgenerator "nokib/campwiz/services/idGenerator"
	rnd "nokib/campwiz/services/round"
	"strings"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ImportService is an interface for importing data from different sources
// All the importer services should implement this interface
type IImportSource interface {
	// This method would be called in a loop to fetch each batch of images
	// It should return the images that were successfully imported and the images that failed to import
	// If there are no images to import it should return nil
	// If there are failed images it should return the reason as a map
	ImportImageResults(failedImageReason *map[string]string) ([]database.ImageResult, *map[string]string)
}
type IDistributionStrategy interface {
	AssignJuries(tx *gorm.DB, round *database.Round, juries []database.Role) (int, error)
}

type TaskRunner struct {
	TaskId               database.IDType
	ImportSource         IImportSource
	DistributionStrategy IDistributionStrategy
}

func NewImportTaskRunner(taskId database.IDType, importService IImportSource) *TaskRunner {
	return &TaskRunner{
		TaskId:       taskId,
		ImportSource: importService,
	}
}
func NewDistributionTaskRunner(taskId database.IDType, strategy IDistributionStrategy) *TaskRunner {
	return &TaskRunner{
		TaskId:               taskId,
		DistributionStrategy: strategy,
	}
}
func (b *TaskRunner) importImagws(conn *gorm.DB, task *database.Task) (successCount, failedCount int) {
	round_repo := database.NewRoundRepository()
	round, err := round_repo.FindByID(conn, *task.AssociatedRoundID)
	if err != nil {
		log.Println("Error fetching round: ", err)
		return
	}
	if round.LatestDistributionTaskID != nil && *round.LatestDistributionTaskID != task.TaskID {
		// log.Println("Task is not the latest task for the round")
		// task.Status = database.TaskStatusFailed
		// return
	}
	currentRoundStatus := round.Status
	{
		// Update the round status to importing
		round.LatestDistributionTaskID = &task.TaskID
		round.Status = database.RoundStatusImporting
		conn.Save(round)
		defer func() {
			round.Status = currentRoundStatus
			conn.Save(round)
		}()
	}
	FailedImages := &map[string]string{}
	technicalJudge := rnd.NewTechnicalJudgeService(round)
	user_repo := database.NewUserRepository()
	for {
		successBatch, failedBatch := b.ImportSource.ImportImageResults(FailedImages)
		if failedBatch != nil {
			task.FailedCount = len(*failedBatch)
			*task.FailedIds = datatypes.NewJSONType(*failedBatch)
		}
		if successBatch == nil {
			break
		}
		images := []database.ImageResult{}
		log.Println("Processing batch of images")
		for _, image := range successBatch {
			if technicalJudge.PreventionReason(image) != "" {
				images = append(images, image)
			}
		}
		successCount += len(images)
		task.SuccessCount = successCount
		participants := map[database.UserName]database.IDType{}
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
	{
		task.Status = database.TaskStatusSuccess
		round.LatestDistributionTaskID = nil // Reset the latest task id
	}
	return
}

func (b *TaskRunner) distributeEvaluations(tx *gorm.DB, task *database.Task) (successCount, failedCount int, err error) {
	round_repo := database.NewRoundRepository()
	round, err := round_repo.FindByID(tx, *task.AssociatedRoundID)
	if err != nil {
		log.Println("Error fetching round: ", err)
		return
	}
	jury_repo := database.NewRoleRepository()
	filter := &database.RoleFilter{
		RoundID:    round.RoundID,
		CampaignID: round.CampaignID,
	}
	j := database.RoleTypeJury
	filter.Type = &j
	// cacheDB, closeCache := cache.GetCacheDB()
	// defer closeCache()
	// _, err = cache.ExportToCache(tx, cacheDB, &database.EvaluationFilter{
	// 	// RoundID:    round.RoundID,
	// 	// CampaignID: round.CampaignID,
	// 	CommonFilter: database.CommonFilter{
	// 		Limit: 50,
	// 	},
	// }) //, task.TaskID)
	// if err != nil {
	// 	log.Println("Error exporting to cache: ", err)
	// 	cacheDB.Rollback()
	// 	return
	// }
	// cacheDB.Commit()
	juries, err := jury_repo.ListAllRoles(tx, filter)
	if err != nil {
		log.Println("Error fetching juries: ", err)
		return
	}
	if len(juries) == 0 {
		log.Println("No juries found")
		return
	}
	log.Printf("Found %d juries\n", len(juries))

	successCount, err = b.DistributionStrategy.AssignJuries(tx, round, juries)
	if err != nil {
		log.Println("Error assigning juries: ", err)
		tx.Rollback()
		return
	}
	return
}

func (b *TaskRunner) Run() {
	task_repo := database.NewTaskRepository()
	conn, close := database.GetDB()
	defer close()

	task, err := task_repo.FindByID(conn, b.TaskId)
	if err != nil {
		log.Println("Error fetching task: ", err)
		return
	}
	defer func() {
		res := conn.Updates(&task)
		if res.Error != nil {
			log.Println("Error saving task: ", res.Error)
		}
	}()
	successCount, failedCount := 0, 0
	if task.Type == database.TaskTypeImportFromCommons {
		if b.ImportSource == nil {
			log.Printf("No import source found for task %s\n", b.TaskId)
			task.Status = database.TaskStatusFailed
			return
		}
		successCount, failedCount = b.importImagws(conn, task)
		task.Status = database.TaskStatusSuccess
	} else if task.Type == database.TaskTypeDistributeEvaluations {
		tx := conn.Begin()
		successCount, failedCount, err = b.distributeEvaluations(tx, task)
		if err != nil {
			log.Println("Error distributing evaluations: ", err)
			tx.Rollback()
			task.Status = database.TaskStatusFailed
			return
		}
		task.Status = database.TaskStatusSuccess
		tx.Commit()
	} else {
		log.Printf("Unknown task type %s\n", task.Type)
		task.Status = database.TaskStatusFailed
		return
	}
	task.Status = database.TaskStatusSuccess
	log.Printf("Task %s completed with %d success and %d failed\n", b.TaskId, successCount, failedCount)
}
