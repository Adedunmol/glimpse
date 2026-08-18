package router

import (
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerDeviceRoutes(r *echo.Echo, h *handler.Handlers, auth *middleware.AuthMiddleware) {
	r.POST("/register", h.Device.RegisterDevice, auth.RequireAuth)
}
