package database

import (
	"nokib/campwiz/consts"
	"time"
)

type User struct {
	UserID       string                 `json:"id" gorm:"primaryKey;type:varchar(255)"`
	RegisteredAt time.Time              `json:"registeredAt"`
	Username     string                 `json:"username" gorm:"unique,not null"`
	Permission   consts.PermissionGroup `json:"permission" gorm:"type:bigint;default:0"`
}
