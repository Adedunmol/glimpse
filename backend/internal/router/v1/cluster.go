package v1

import (
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/labstack/echo/v4"
)

func registerClusterRoutes(r *echo.Group, h *handler.ClusterHandler, auth *middleware.AuthMiddleware) {
	clusterGroup := r.Group("/clusters")
	clusterGroup.Use(auth.RequireAuth)

	clusterGroup.GET("", h.GetClusters)
	dynamic := clusterGroup.Group("/:clusterId")
	dynamic.GET("", h.GetClusterID)
	dynamic.POST("/links", h.CreateLink)
}
