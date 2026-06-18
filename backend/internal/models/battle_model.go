package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type TroopDropBody struct {
	Name string `json:"name"`
	X uint8 `json:"x"`
	Y uint8 `json:"y"`
}
type BattleManager struct {
    Mu      sync.Mutex
    Battles map[string][]*Client
}
type Client struct {
    ID       string
    Conn     *websocket.Conn
    Send     chan []byte   // outgoing
    Incoming chan []byte   // incoming
}

func (c *Client) Write() {
	for data := range c.Send {
		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			data,
		)

		if err != nil {
			return
		}
	}
}

func (c *Client) Read() {
    defer close(c.Incoming)

	for {
		messageType, message, err := c.Conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType != websocket.TextMessage {
			continue
		}

		c.Incoming <- message
	}
}