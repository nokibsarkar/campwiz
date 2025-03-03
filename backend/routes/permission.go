package routes

import (
	"fmt"
	"nokib/campwiz/consts"
	"nokib/campwiz/database/cache"

	"github.com/gin-gonic/gin"
)

func WithPermission(requiredPermission consts.Permission, handler func(c *gin.Context, sess *cache.Session)) gin.HandlerFunc {
	return WithSession(func(c *gin.Context, sess *cache.Session) {
		fmt.Println("Checking permission: ", sess.Permission, requiredPermission)
		if sess.Permission.HasPermission(requiredPermission) {
			handler(c, sess)
		} else {
			c.JSON(403, ResponseError{Detail: "Permission denied"})
		}
	})
}
