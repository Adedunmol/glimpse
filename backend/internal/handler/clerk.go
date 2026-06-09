package handler

import (
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/labstack/echo/v4"
)

type ClerkWebHookHandler struct {
	Handler
}

func NewClerkWebHookHandler(s *server.Server) *ClerkWebHookHandler {
	return &ClerkWebHookHandler{
		Handler: NewHandler(s),
	}
}

func (clerk *ClerkWebHookHandler) HandleEvent(c echo.Context) error {

	return nil
}
