package ws

import (
	"log"

	"github.com/hertz-contrib/websocket"
)

type Client struct {
	UserID   string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}

func (c *Client) WriteMessage() {
	defer c.Conn.Close()

	for msg := range c.Send {
		if err := c.Conn.WriteMessage(
			websocket.TextMessage, msg,
		); err != nil {
			log.Printf("Failed to send message to user %s: %v", c.UserID, err)
		}
	}
}

func (c *Client) ReadMessage() {
	defer func() {
		ManagerInstance.RemoveClient(c.UserID)
		close(c.Send)
	}()

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message from user %s: %v", c.UserID, err)
			break
		}

		ManagerInstance.BoadcastMessage(msg, c)
	}
}
