package services

import (
	"fmt"
	"nokib/campwiz/consts"
	"nokib/campwiz/database"

	"gorm.io/gorm"
	"k8s.io/apimachinery/pkg/util/sets"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetUserByID(conn *gorm.DB, id database.IDType) (*database.User, error) {
	userFilter := &database.User{UserID: id}
	user := &database.User{}
	result := conn.First(user, userFilter)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		return nil, result.Error
	}
	return user, nil

}
func (u *UserService) GetUserByUsername(conn *gorm.DB, username string) (*database.User, error) {
	userFilter := &database.User{Username: username}
	user := &database.User{}
	result := conn.First(user, userFilter)
	if result.Error != nil {
		fmt.Println("Error: ", result.Error)
		return nil, result.Error
	}
	return user, nil
}
func (u *UserService) GetOrCreateUser(conn *gorm.DB, user *database.User) (*database.User, error) {
	result := conn.FirstOrCreate(user, user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}
func (u *UserService) EnsureExists(tx *gorm.DB, usernameSet sets.Set[string]) (map[string]database.IDType, error) {
	user_repo := database.NewUserRepository()
	userName2Id, err := user_repo.FetchExistingUsernames(tx, usernameSet.UnsortedList())
	if err != nil {
		return nil, err
	}
	if len(userName2Id) > 0 {
		for username, _ := range userName2Id {
			usernameSet.Delete(username)
		}
	}
	nonExistentUsers := usernameSet.UnsortedList()
	if len(nonExistentUsers) == 0 {
		return userName2Id, nil
	}
	commons_repo := database.NewCommonsRepository()
	users, err := commons_repo.GeUsersFromUsernames(nonExistentUsers)
	if err != nil {
		return nil, err
	}
	new_users := []database.User{}
	for _, u := range users {
		new_user := database.User{
			UserID:       GenerateID("user"),
			RegisteredAt: u.Registered,
			Username:     u.Name,
			Permission:   consts.PermissionGroupUSER,
		}
		new_users = append(new_users, new_user)
		userName2Id[new_user.Username] = new_user.UserID
	}
	result := tx.Create(new_users)
	return userName2Id, result.Error
}
