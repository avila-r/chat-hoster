package chat

import (
	"log"
	"sync"
)

type Server struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
	mutex      sync.Mutex
}

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}

func NewServer() *Server {
	return &Server{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

// Starts the server loop to handle registering, unregistering, and broadcasting.
func (s *Server) Run() {
	for {
		select {
		case client := <-s.Register:
			s.register(client)

		case client := <-s.Unregister:
			s.unregister(client)

		case message := <-s.Broadcast:
			s.broadcast(message)
		}
	}
}

// Adds a new client to the corresponding room.
func (s *Server) register(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	room, exists := s.Rooms[client.RoomID]

	if !exists {
		log.Printf("Room %s does not exist", client.RoomID)
		return
	}

	if _, exists := room.Clients[client.ID]; !exists {
		room.Clients[client.ID] = client
		log.Printf("Client %s joined room %s", client.Username, room.Name)
	}
}

// Removes a client from the room and broadcasts a message.
func (s *Server) unregister(client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	room, exists := s.Rooms[client.RoomID]

	if !exists {
		log.Printf("Room %s does not exist", client.RoomID)
		return
	}

	if _, exists := room.Clients[client.ID]; exists {
		delete(room.Clients, client.ID)
		close(client.Message)

		log.Printf("Client %s left room %s", client.Username, room.Name)

		if len(room.Clients) > 0 {
			s.Broadcast <- &Message{
				Content:  "User left the chat",
				RoomID:   client.RoomID,
				Username: client.Username,
			}
		} else {
			log.Printf("Room %s is now empty, deleting room", room.ID)
			delete(s.Rooms, room.ID)
		}
	}
}

// broadcast sends a message to all clients in a room.
func (s *Server) broadcast(message *Message) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	room, exists := s.Rooms[message.RoomID]

	if !exists {
		log.Printf("Room %s does not exist", message.RoomID)
		return
	}

	for _, client := range room.Clients {
		select {

		case client.Message <- message:
			log.Printf("Broadcast message to client %s: %s", client.Username, message.Content)

		default:
			log.Printf("Client %s message buffer full, dropping message", client.Username)

		}
	}
}
