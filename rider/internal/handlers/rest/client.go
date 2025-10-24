package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Binit-Dhakal/Saarathi/rider/internal/application"
)

type Client struct {
	id        string
	updateSvc application.RiderUpdateService
	send      chan any
	done      chan struct{}
	flusher   http.Flusher
	w         http.ResponseWriter
}

func (c *Client) writePump() {

	for {
		select {
		case message := <-c.send:
			payload, err := json.Marshal(message)
			if err != nil {
				continue
			}
			if _, err := fmt.Fprintf(c.w, "data: %s\n\n", payload); err != nil {
				return
			}
			c.flusher.Flush()

		case <-c.done:
			return
		}
	}

}
