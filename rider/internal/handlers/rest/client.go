package rest

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/Binit-Dhakal/Saarathi/pkg/ddd"
	"github.com/Binit-Dhakal/Saarathi/rider/internal/application"
)

type Client struct {
	id         string
	updateSvc  application.RiderUpdateService
	send       chan any
	done       chan struct{}
	flusher    http.Flusher
	w          http.ResponseWriter
	reqContext context.Context
}

func (c *Client) Start() {
	c.done = make(chan struct{})
	go c.writePump()
	go c.readPump()
}

// for SSE; as we dont to read anything just listen for connection closure
func (c *Client) readPump() {
	defer close(c.done)
	<-c.reqContext.Done()
}

func (c *Client) writePump() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message := <-c.send:
			fmt.Fprintf(c.w, "data: %s\n\n", message)
			c.flusher.Flush()

		case <-ticker.C:
			fmt.Fprintf(c.w, ": keep-alive\n\n")
			c.flusher.Flush()

		case <-c.done:
			return
		}
	}

}
