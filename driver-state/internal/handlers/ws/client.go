package ws

import (
	"encoding/json"
	"log"
	"time"

	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/application"
	"github.com/Binit-Dhakal/Saarathi/driver-state/internal/dto"
	"github.com/gorilla/websocket"
)

type Client struct {
	ID          string
	Conn        *websocket.Conn
	Send        chan any
	locationSvc application.LocationService
	done        chan struct{}
}

func (c *Client) Start() {
	c.done = make(chan struct{})
	go c.readPump()
	go c.writePump()
}

func (c *Client) readPump() {
	defer c.cleanup()

	for {
		var baseMessage dto.BaseMessage
		err := c.Conn.ReadJSON(&baseMessage)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Driver %s closed with an error: %v\n", c.ID, err)
			} else {
				log.Printf("Failed to read message: %v\n", err)
			}
			break
		}

		switch baseMessage.Event {
		case "DRIVER_LOCATION_UPDATE":
			var locationPayload dto.DriverLocationMessage
			if err := json.Unmarshal(baseMessage.Data, &locationPayload); err != nil {
				log.Println("Failed to unmarshal location payload: ", err)
				continue
			}
			// How to do this step
			err = c.locationSvc.UpsertDriverLocation(&locationPayload, c.ID)
			if err != nil {
				log.Println("Failed to save drivers location: ", err)
				continue
			}
		}
	}
}

func (c *Client) writePump() {
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("Failed to marshal message for driver %s:%v\n", c.ID, err)
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Printf("Failed to write message to driver: %s:%v\n", c.ID, err)
				return
			}

		case <-c.done:
			return
		}
	}
}

func (c *Client) cleanup() {
	select {
	case <-c.done:
		return
	default:
		close(c.done)
		close(c.Send)
		c.Conn.Close()
		log.Printf("Driver %s fully cleaned up", c.ID)
	}
}
