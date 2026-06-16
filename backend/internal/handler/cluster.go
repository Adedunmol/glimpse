package handler

import (
	"net/http"

	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/cluster"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
	"github.com/labstack/echo/v4"
)

type ClusterHandler struct {
	Handler
	clusterService *service.ClusterService
}

func NewClusterHandler(s *server.Server, clusterService *service.ClusterService) *ClusterHandler {
	return &ClusterHandler{
		Handler:        NewHandler(s),
		clusterService: clusterService,
	}
}

func (h *ClusterHandler) GetClusterID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *cluster.GetClusterByIDPayload) (*cluster.Cluster, error) {
			userID := middleware.GetUserID(c)
			logger := middleware.GetLogger(c)

			return h.clusterService.GetClusterByID(c.Request().Context(), logger, userID, payload.ID)
		},
		http.StatusOK,
		&cluster.GetClusterByIDPayload{},
	)(c)
}

func (h *ClusterHandler) GetClusters(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *cluster.GetClustersQuery) (*model.PaginatedResponse[cluster.Cluster], error) {
			userID := middleware.GetUserID(c)
			logger := middleware.GetLogger(c)
			return h.clusterService.GetClusters(c.Request().Context(), logger, userID, query)
		},
		http.StatusOK,
		&cluster.GetClustersQuery{},
	)(c)
}
