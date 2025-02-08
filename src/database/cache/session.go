package cache

import (
	"nokib/campwiz/consts"
	"time"
)

type Session struct {
	ID         string                 `json:"id" gorm:"primaryKey"`
	UserID     string                 `json:"userId"`
	Username   string                 `json:"username"`
	Permission consts.PermissionGroup `json:"permission"`
	ExpiresAt  time.Time              `json:"expiresAt"`
}
