package gincontext

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var Clients = make(map[string]chan string)

// SSE处理器
func SSEHandler(c *gin.Context) {
	uid := c.Param("uid")
	messageChan := make(chan string)
	Clients[uid] = messageChan
	defer delete(Clients, uid)

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		select {
		case message := <-messageChan:
			fmt.Fprintf(c.Writer, "data: %s\n\n", message)
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
