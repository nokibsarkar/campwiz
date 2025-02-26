package database

import (
	"log"
	"nokib/campwiz/consts"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID       IDType                 `json:"id" gorm:"primaryKey"`
	RegisteredAt time.Time              `json:"registeredAt"`
	Username     string                 `json:"username" gorm:"unique;not null;index"`
	Permission   consts.PermissionGroup `json:"permission" gorm:"type:bigint;default:0"`
}

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}
func (u *UserRepository) FetchExistingUsernames(conn *gorm.DB, usernames []string) (map[string]IDType, error) {
	type APIUser struct {
		Username string
		UserID   IDType
	}
	exists := []APIUser{}

	res := conn.Model(&User{}).Limit(len(usernames)).Find(&APIUser{}, "username IN (?)", usernames).Find(&exists)
	if res.Error != nil {
		return nil, res.Error
	}
	userName2IDMap := map[string]IDType{}
	for _, u := range exists {
		userName2IDMap[u.Username] = u.UserID
	}
	return userName2IDMap, nil

}
func (u *UserRepository) EnsureExists(tx *gorm.DB, usernameToRandomIdMap map[string]IDType) (map[string]IDType, error) {
	usernames := []string{}
	for username := range usernameToRandomIdMap {
		usernames = append(usernames, username)
	}
	userName2Id, err := u.FetchExistingUsernames(tx, usernames)
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
	log.Println("Non existent users: ", nonExistentUsers)
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
	if len(new_users) == 0 {
		return userName2Id, nil
	}
	result := tx.Create(new_users)
	return userName2Id, result.Error
}
