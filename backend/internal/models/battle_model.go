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
    Mu      *sync.Mutex
    Battles map[string][]*Client
}
type Client struct {
    ID       string
    Conn     *websocket.Conn
    Send     chan []byte   // outgoing
    Incoming chan []byte   // incoming
	Done     chan struct{}
}

func (c *Client) Write() {
	for {
		select {
		case data := <- c.Send:
		err := c.Conn.WriteMessage(
			websocket.TextMessage,
			data,
		)
		if err != nil {
			return
		}

		case <- c.Done: 
			return 
		}
	}
}

func (c *Client) Read() {
    defer close(c.Incoming)
	defer close(c.Done)

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