package router

import (
	"github.com/avila-r/chat-hoster/db"
	"github.com/avila-r/chat-hoster/internal/chat"
	"github.com/avila-r/chat-hoster/internal/users"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

var (
	user_repository users.Repository
	user_service    users.Service
	user_handler    users.Handler
	server          chat.Server
	chat_handler    chat.Handler
)

func setup() {
	user_repository = *users.NewRepository(db.GetConnection())
	user_service = *users.NewService(user_repository)
	user_handler = *users.NewHandler(user_service)

	server = *chat.NewServer()
	chat_handler = *chat.NewHandler(&server)
	go server.Run()
}

func EnableRouting(app *fiber.App) {
	setup()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Post(
		"/register",
		user_handler.CreateUser,
	)

	app.Post(
		"/login",
		user_handler.Login,
	)

	app.Post(
		"/logout",
		user_handler.Logout,
	)

	app.Get(
		"/clients/:room_id",
		chat_handler.GetClients,
	)

	app.Get(
		"/rooms",
		chat_handler.GetRooms,
	)

	app.Post(
		"/rooms/:room_id",
		chat_handler.JoinRoom,
	)

	app.Post(
		"/rooms",
		chat_handler.CreateRoom,
	)

	_ = app.Listen(":8888")
}
