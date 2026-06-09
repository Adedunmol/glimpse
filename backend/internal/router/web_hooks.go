package router

import (
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/labstack/echo/v4"
)

func registerWebHookRoutes(r *echo.Echo, h *handler.Handlers) {
	r.POST("/clerk/webhook", nil)
}
