package chat

type Server struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
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

func (s *Server) Run() {
	for {
		select {

		case channel := <-s.Register:
			room, ok := s.Rooms[channel.RoomID]

			if ok {
				_, ok := room.Clients[channel.ID]

				if !ok {
					room.Clients[channel.ID] = channel
				}
			}

		case channel := <-s.Unregister:
			room, ok := s.Rooms[channel.RoomID]

			if ok {
				_, ok := room.Clients[channel.ID]

				if ok {
					if len(room.Clients) != 0 {
						s.Broadcast <- &Message{
							Content:  "user left the chat",
							RoomID:   channel.RoomID,
							Username: channel.Username,
						}
					}

					delete(room.Clients, channel.ID)
					close(channel.Message)
				}
			}

		case message := <-s.Broadcast:
			room, ok := s.Rooms[message.RoomID]

			if ok {
				for _, client := range room.Clients {
					client.Message <- message
				}
			}
		}
	}
}
