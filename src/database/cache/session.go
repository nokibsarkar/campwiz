package cache

import (
	"nokib/campwiz/consts"
	"nokib/campwiz/database"
	"time"
)

type Session struct {
	ID         database.IDType        `json:"id" gorm:"primaryKey"`
	UserID     database.IDType        `json:"userId"`
	Username   database.UserName      `json:"username"`
	Permission consts.PermissionGroup `json:"permission"`
	ExpiresAt  time.Time              `json:"expiresAt"`
}
