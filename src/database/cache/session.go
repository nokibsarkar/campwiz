package cache

import (
	"nokib/campwiz/consts"
	"time"
)

type Session struct {
	ID         uint64                 `json:"id" gorm:"primaryKey,autoIncrement"`
	UserID     string                 `json:"userId"`
	Username   string                 `json:"username"`
	Permission consts.PermissionGroup `json:"permission"`
	ExpiresAt  time.Time              `json:"expiresAt"`
}
