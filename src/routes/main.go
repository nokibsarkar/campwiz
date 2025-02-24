package routes

import (
	"fmt"

	"nokib/campwiz/database/cache"

	"github.com/gin-gonic/gin"
)

func WithSession(callback func(*gin.Context, *cache.Session)) gin.HandlerFunc {
	return func(c *gin.Context) {
		session := GetSession(c)
		if session == nil {
			c.JSON(401, ResponseError{
				Detail: "Internal Server Error : Session not found",
			})
			return
		}
		fmt.Println("Session: ", session)
		callback(c, session)
	}
}
func GetSession(c *gin.Context) *cache.Session {
	sess, ok := c.Get(SESSION_KEY)
	if !ok {
		return nil
	}
	session, ok := sess.(*cache.Session)
	if !ok {
		return nil
	}
	return session
}
func NewRoutes(nonAPIParent *gin.RouterGroup) {
	r := nonAPIParent.Group("/api/v2")
	authenticatorService := NewAuthenticationService()
	r.Use(authenticatorService.Authenticate)
	NewStaticRouter(nonAPIParent)
	NewUserAuthenticationRoutes(nonAPIParent)
	NewCampaignRoutes(r)
	NewSubmissionRoutes(r)
	NewUserRoutes(r)
	NewRoundRoutes(r)
	NewBatchRouter(r)
}
