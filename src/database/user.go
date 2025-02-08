package database

import (
	"nokib/campwiz/consts"
	"time"
)

type User struct {
	ID           string                 `json:"id" gorm:"primaryKey"`
	RegisteredAt time.Time              `json:"registeredAt"`
	Username     string                 `json:"username" gorm:"unique,not null"`
	Permission   consts.PermissionGroup `json:"permission" gorm:"type:bigint;default:0"`
}
