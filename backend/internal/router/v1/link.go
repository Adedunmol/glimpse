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

	// lookups by link id and token live under their own prefixes: a bare
	// /links/:param can only ever bind one name, so the three used to collide
	// and Echo resolved every request to the last one registered
	linkGroup.GET("/id/:linkId", h.GetLinkByID)
	linkGroup.GET("/token/:token", h.GetLinkByToken)

	linkGroup.GET("/:clusterId", h.GetUploadsByClusterID)
}
