package routes

import (
	"fmt"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

func SSEHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")
		c.Next()
	}
}
func SseHandler(c *gin.Context) {
	// ...
	c.Stream(func(w io.Writer) bool {
		time.Sleep(time.Second * 1)
		now := time.Now().Format("2006-01-02 15:04:05")
		currentTime := fmt.Sprintf("The Current Time Is %v", now)
		c.SSEvent("message", currentTime)
		c.Writer.Flush()
		fmt.Println("Sent: ", currentTime)
		return true
	})
}
func NewSseRoutes(parent *gin.RouterGroup) {
	r := parent.Group("/sse")
	r.Use(SSEHeadersMiddleware())
	r.GET("/sse", SseHandler)
}
