package gincontext

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

var Clients = make(map[string]chan string)

// SSEHandler SSE推送
func SSEHandler(c *gin.Context) {
	uid := c.Param("uid")
	messageChan := make(chan string)
	Clients[uid] = messageChan
	defer delete(Clients, uid)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(30 * time.Second) // 30秒发送一次心跳
	defer ticker.Stop()

	for {
		select {
		case message := <-messageChan:
			fmt.Fprintf(c.Writer, "data: %s\n\n", message)
			c.Writer.Flush()
		case <-ticker.C:
			fmt.Fprintf(c.Writer, ": keep-alive\n\n")
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
