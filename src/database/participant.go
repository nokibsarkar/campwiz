package database

import (
	"nokib/campwiz/consts"

	"gorm.io/gorm"
)

type Participant struct {
	ParticipantID IDType `json:"id" gorm:"primaryKey"`
	Username      string `json:"name"`
}
type ParticipantRepository struct{}

func NewParticipantRepository() *ParticipantRepository {
	return &ParticipantRepository{}
}
func (u *ParticipantRepository) FetchExistingUsernames(conn *gorm.DB, usernames []string) (map[string]string, error) {
	exists := []Participant{}

	res := conn.Model(&Participant{}).Limit(len(usernames)).Find("username IN (?)", usernames).Find(&exists)
	if res.Error != nil {
		return nil, res.Error
	}
	userName2IDMap := map[string]string{}
	for _, u := range exists {
		userName2IDMap[u.Username] = string(u.ParticipantID)
	}
	return userName2IDMap, nil
}
func (u *ParticipantRepository) EnsureExists(tx *gorm.DB, usernameToRandomIdMap map[string]IDType) (map[string]IDType, error) {
	usernames := []string{}
	for username := range usernameToRandomIdMap {
		usernames = append(usernames, username)
	}
	user_repo := NewUserRepository()
	userName2Id, err := user_repo.FetchExistingUsernames(tx, usernames)
	if err != nil {
		return nil, err
	}
	if len(userName2Id) > 0 {
		for username := range userName2Id {
			delete(usernameToRandomIdMap, username)
		}
	}
	if len(usernameToRandomIdMap) == 0 {
		return userName2Id, nil
	}
	nonExistentUsers := make([]string, 0, len(usernameToRandomIdMap))
	for nonExistingUsername := range usernameToRandomIdMap {
		nonExistentUsers = append(nonExistentUsers, nonExistingUsername)
	}
	commons_repo := NewCommonsRepository()
	users, err := commons_repo.GeUsersFromUsernames(nonExistentUsers)
	if err != nil {
		return nil, err
	}
	new_users := []User{}
	for _, u := range users {
		new_user := User{
			UserID:       usernameToRandomIdMap[u.Name],
			RegisteredAt: u.Registered,
			Username:     u.Name,
			Permission:   consts.PermissionGroupUSER,
		}
		new_users = append(new_users, new_user)
		userName2Id[new_user.Username] = new_user.UserID
	}
	result := tx.Create(new_users)
	if result.Error != nil {
		return nil, result.Error
	}
	new_participants := []Participant{}
	for _, p := range new_users {
		new_participants = append(new_participants, Participant{
			ParticipantID: IDType(p.UserID),
			Username:      p.Username,
		})
		userName2Id[p.Username] = p.UserID
	}
	result = tx.Create(new_participants)
	return userName2Id, result.Error
}
