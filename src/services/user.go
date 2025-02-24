package services

import (
	"fmt"
	"nokib/campwiz/database"

	"gorm.io/gorm"
)

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (u *UserService) GetUserByID(conn *gorm.DB, id string) (*database.User, error) {
	userFilter := &database.User{UserID: id}
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
