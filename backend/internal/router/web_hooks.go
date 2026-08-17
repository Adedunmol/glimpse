package router

import (
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerWebHookRoutes(r *echo.Echo, h *handler.Handlers) {
	// the raw body is needed for signature verification, so capture it before
	// the binder consumes it
	r.POST("/clerk/webhook", h.Clerk.HandleEvent, middleware.CaptureRawBody())
}
