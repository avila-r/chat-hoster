package chat

import (
	http "github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type Handler struct {
	server *Server
}

func NewHandler(s *Server) *Handler {
	return &Handler{s}
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
			"error": err.Error(),
		})
	}

	h.server.Rooms[request.ID] = &Room{
		ID:      request.ID,
		Name:    request.Name,
		Clients: make(map[string]*Client),
	}

	return c.JSON(http.StatusCreated)
}

func (h *Handler) JoinRoom(c *http.Ctx) error {
	if !websocket.IsWebSocketUpgrade(c) {
		return c.Status(http.StatusUpgradeRequired).JSON(http.Map{
			"error": "This endpoint requires a WebSocket connection",
		})
	}

	return websocket.New(func(connection *websocket.Conn) {
		roomID := c.Params("roomId")
		clientID := c.Query("userId")
		username := c.Query("username")

		client := &Client{
			Connection: connection,
			Message:    make(chan *Message, 10),
			ID:         clientID,
			RoomID:     roomID,
			Username:   username,
		}

		m := &Message{
			Content:  "A new user has joined the room",
			RoomID:   roomID,
			Username: username,
		}

		h.server.Register <- client
		h.server.Broadcast <- m

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
		rooms []RoomResponse
	)

	return c.Status(http.StatusOK).JSON(rooms)
}

func (h *Handler) GetClients(c *http.Ctx) error {
	type ClientResponse struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}

	var (
		clients []ClientResponse
	)

	return c.Status(http.StatusOK).JSON(clients)
}
