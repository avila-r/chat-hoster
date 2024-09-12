package chat

import (
	"github.com/gofiber/websocket/v2"
)

type Client struct {
	Connection *websocket.Conn
	Message    chan *Message
	ID         string `json:"id"`
	RoomID     string `json:"roomId"`
	Username   string `json:"username"`
}

type Message struct {
	RoomID   string `json:"roomId"`
	Username string `json:"username"`
	Content  string `json:"content"`
}

func (c *Client) Write() {
	defer func() {
		c.Connection.Close()
	}()

	for {
		message, ok := <-c.Message

		if !ok {
			return
		}

		c.Connection.WriteJSON(message)
	}
}

func (c *Client) Read(s *Server) {
	defer func() {
		s.Unregister <- c

		c.Connection.Close()
	}()

	for {
		_, content, err := c.Connection.ReadMessage()

		if err != nil {
			break
		}

		message := &Message{
			Content:  string(content),
			RoomID:   c.RoomID,
			Username: c.Username,
		}

		s.Broadcast <- message
	}
}
