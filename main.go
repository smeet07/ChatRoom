package main

import (
	"github.com/Smeet07/ChatAppGo/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
)

func main() {
	viewsEngine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: viewsEngine,
	})
	app.Static("/static", "./static")
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to fiber")
	})
	appHandler := handlers.NewAppHandler()
	// Add appHandler routes
	app.Get("/", appHandler.HandleGetIndex)
	app.Post("/join", appHandler.HandleJoinRoom)
	app.Get("/room/:code", appHandler.HandleGetRoom)
	server := NewWebSocket()
	app.Get("/ws/:code", websocket.New(func(ctx *websocket.Conn) {
		server.HandleWebSocket(ctx, ctx.Params("code"))
	}))

	go server.HandleMessages()
	app.Listen(":3000")

}
