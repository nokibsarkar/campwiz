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
	sess, ok := c.Get("session")
	if !ok {
		return nil
	}
	session, ok := sess.(*cache.Session)
	return session
}
func NewRoutes(parent *gin.RouterGroup) {

	UserAuthenticationRoutes(parent)
	r := parent.Group("/api/v2")
	authenticatorService := NewAuthenticationService()
	r.Use(authenticatorService.Authenticate)
	NewCampaignRoutes(r)
	NewSubmissionRoutes(r)
	NewUserRoutes(r)
}
