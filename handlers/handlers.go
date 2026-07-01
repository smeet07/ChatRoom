package handlers

import (
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type AppHandler struct{}

func NewAppHandler() *AppHandler {
	return &AppHandler{}
}
func (a *AppHandler) HandleGetIndex(ctx *fiber.Ctx) error {
	context := fiber.Map{}
	return ctx.Render("index", context)
}

func (a *AppHandler) HandleJoinRoom(ctx *fiber.Ctx) error {
	user := strings.TrimSpace(ctx.FormValue("user"))
	roomCode := strings.TrimSpace(ctx.FormValue("room"))
	if user == "" || roomCode == "" {
		return ctx.Redirect("/")
	}

	return ctx.Redirect("/room/" + url.PathEscape(roomCode) + "?user=" + url.QueryEscape(user))
}

func (a *AppHandler) HandleGetRoom(ctx *fiber.Ctx) error {
	roomCode := strings.TrimSpace(ctx.Params("code"))
	user := strings.TrimSpace(ctx.Query("user"))
	if roomCode == "" || user == "" {
		return ctx.Redirect("/")
	}

	return ctx.Render("room", fiber.Map{
		"RoomCode": roomCode,
		"User":     user,
	})
}
