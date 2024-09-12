package chat

import (
	"log"

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

// Write listens on the Message channel and sends messages to the WebSocket connection.
func (c *Client) Write() {
	defer func() {
		c.Connection.Close()
		log.Printf("Connection closed for client: %s", c.ID)
	}()

	for message := range c.Message {
		if err := c.Connection.WriteJSON(message); err != nil {
			log.Printf("Error sending message to client %s: %v", c.ID, err)
			return
		}
	}

	log.Printf("Message channel closed for client: %s", c.ID)
}

// Read reads messages from the WebSocket connection and sends them to the server's broadcast channel.
func (c *Client) Read(s *Server) {
	defer func() {
		s.Unregister <- c
		c.Connection.Close()
		log.Printf("Client %s unregistered and connection closed", c.ID)
	}()

	for {
		_, content, err := c.Connection.ReadMessage()

		if err != nil {
			log.Printf("Error reading message from client %s: %v", c.ID, err)
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
