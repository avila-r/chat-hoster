package chat

import (
	"log"
	"sync"

	http "github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handler struct {
	server *Server
	mutex  sync.RWMutex
}

func NewHandler(s *Server) *Handler {
	return &Handler{
		server: s,
	}
}

func (h *Handler) CreateRoom(c *http.Ctx) error {
	type CreateRoomRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var (
		request CreateRoomRequest
	)

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(http.Map{
			"error": "failed to parse request body - " + err.Error(),
		})
	}

	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, exists := h.server.Rooms[request.ID]; exists {
		return c.Status(http.StatusConflict).JSON(http.Map{
			"error": "room already exists",
		})
	}

	h.server.Rooms[request.ID] = &Room{
		ID:      request.ID,
		Name:    request.Name,
		Clients: make(map[string]*Client),
	}

	log.Printf("Room created: %s", request.Name)

	return c.Status(http.StatusCreated).JSON(http.Map{
		"id":   request.ID,
		"name": request.Name,
	})
}

func (h *Handler) JoinRoom(c *http.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return c.Status(http.StatusUpgradeRequired).JSON(http.Map{
			"error": "this endpoint requires a web socket connection",
		})
	}

	return websocket.New(func(connection *websocket.Conn) {
		room_id := c.Params("room_id")
		client_id := c.Query("user_id")
		username := c.Query("username")

		client := &Client{
			Connection: connection,
			Message:    make(chan *Message, 10),
			ID:         client_id,
			RoomID:     room_id,
			Username:   username,
		}

		m := &Message{
			Content:  "A new user has joined the room",
			RoomID:   room_id,
			Username: username,
		}

		h.server.Register <- client
		h.server.Broadcast <- m

		log.Printf(
			"User %s joined room %s",
			username, room_id,
		)

		go client.Write()
		client.Read(h.server)
	})(c)
}

func (h *Handler) GetRooms(c *http.Ctx) error {
	type RoomResponse struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}

	var (
		rooms = make([]RoomResponse, 0)
	)

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for _, r := range h.server.Rooms {
		rooms = append(rooms, RoomResponse{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	return c.Status(http.StatusOK).JSON(rooms)
}

func (h *Handler) GetClients(c *http.Ctx) error {
	type ClientResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	var (
		clients = make([]ClientResponse, 0)
	)

	h.mutex.RLock()
	defer h.mutex.RUnlock()

	room, ok := h.server.Rooms[c.Params("room_id")]

	if !ok {
		return c.Status(http.StatusNotFound).JSON(http.Map{
			"error": "room not found",
		})
	}

	for _, c := range room.Clients {
		clients = append(clients, ClientResponse{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	return c.Status(http.StatusOK).JSON(clients)
}
