package services

import (
	"errors"
	"nokib/campwiz/database"
)

type EvaluationService struct{}

func NewEvaluationService() *EvaluationService {
	return &EvaluationService{}
}

type EvaluationRequest struct {
	VoteScore    *int   `json:"voteScore"`
	Comment      string `json:"comment"`
	VotePassed   *bool  `json:"votePassed"`
	VotePosition *int   `json:"votePosition"`
}

func (e *EvaluationService) Evaluate(currentUserID database.IDType, evaluationID database.IDType, evaluationRequest *EvaluationRequest) (*database.Evaluation, error) {
	ev_repo := database.NewEvaluationRepository()
	user_repo := database.NewUserRepository()
	jury_repo := database.NewJuryRepository()
	conn, close := database.GetDB()
	defer close()
	tx := conn.Begin()
	evaluation, err := ev_repo.FindEvaluationByID(tx, evaluationID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if evaluation == nil {
		tx.Rollback()
		return nil, errors.New("evaluation not found")
	}
	if evaluation.Type == database.EvaluationTypeBinary && evaluationRequest.VotePassed == nil {
		tx.Rollback()
		return nil, errors.New("votePassed is required for binary evaluation")
	} else if evaluation.Type == database.EvaluationTypeRanking && evaluationRequest.VotePosition == nil {
		tx.Rollback()
		return nil, errors.New("votePosition is required for positional evaluation")
	} else if evaluation.Type == database.EvaluationTypeScore && evaluationRequest.VoteScore == nil {
		tx.Rollback()
		return nil, errors.New("voteScore is required for score evaluation")
	}

	submission := evaluation.Submission
	currentUser, err := user_repo.FindByID(tx, currentUserID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if currentUser == nil {
		tx.Rollback()
		return nil, errors.New("user not found")
	}
	if submission.SubmittedByID == currentUser.UserID {
		tx.Rollback()
		return nil, errors.New("user can't evaluate his/her own submission")
	}
	round := submission.CurrentRound
	campaign := round.Campaign
	juries, err := jury_repo.ListAllRoles(tx, &database.RoleFilter{RoundID: round.RoundID, CampaignID: campaign.CampaignID})
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	juryMap := map[database.IDType]*database.Role{}
	for _, jury := range juries {
		juryMap[jury.RoleID] = &jury
	}
	if _, ok := juryMap[currentUser.UserID]; !ok {
		tx.Rollback()
		return nil, errors.New("user is not a jury")
	}
	return nil, nil
}
func (e *EvaluationService) GetEvaluationById() {
}

func (e *EvaluationService) ListEvaluations(filter *database.EvaluationFilter) ([]database.Evaluation, error) {
	ev_repo := database.NewEvaluationRepository()
	conn, close := database.GetDB()
	defer close()
	return ev_repo.ListAllEvaluations(conn, filter)
}
