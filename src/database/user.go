package database

import (
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
