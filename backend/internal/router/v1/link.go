package v1

import (
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerLinkRoutes(r *echo.Group, h *handler.LinkHandler, auth *middleware.AuthMiddleware) {
	linkGroup := r.Group("/links")
	linkGroup.Use(auth.RequireAuth)

	linkGroup.GET("", h.GetLinks)
	linkGroup.GET("/:linkId", h.GetLinkByID)
	linkGroup.GET("/:token", h.GetLinkByToken)
	linkGroup.GET("/:clusterId", h.GetUploadsByClusterID)
}
