package handler

import (
	"net/http"

	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/link"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
	"github.com/labstack/echo/v4"
)

type LinkHandler struct {
	Handler
	linkService *service.LinkService
}

func NewLinkHandler(s *server.Server, linkService *service.LinkService) *LinkHandler {
	return &LinkHandler{
		Handler:     NewHandler(s),
		linkService: linkService,
	}
}

func (h *LinkHandler) GetLinkByID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *link.GetLinkByIDPayload) (*link.Link, error) {
			userID := middleware.GetUserID(c)
			logger := middleware.GetLogger(c)

			return h.linkService.GetLinkByID(c.Request().Context(), logger, userID, payload.ID)
		},
		http.StatusOK,
		&link.GetLinkByIDPayload{},
	)(c)
}

func (h *LinkHandler) GetLinkByToken(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *link.GetLinkByTokenPayload) (*link.Link, error) {
			userID := middleware.GetUserID(c)
			logger := middleware.GetLogger(c)

			return h.linkService.GetLinkByToken(c.Request().Context(), logger, userID, payload.Token)
		},
		http.StatusOK,
		&link.GetLinkByTokenPayload{},
	)(c)
}

func (h *LinkHandler) GetLinks(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, query *link.GetLinksQuery) (*model.PaginatedResponse[link.Link], error) {
			userID := middleware.GetUserID(c)
			logger := middleware.GetLogger(c)
			return h.linkService.GetLinks(c.Request().Context(), logger, userID, query)
		},
		http.StatusOK,
		&link.GetLinksQuery{},
	)(c)
}

func (h *LinkHandler) GetUploadsByClusterID(c echo.Context) error {
	return Handle(
		h.Handler,
		func(c echo.Context, payload *link.GetLinkByClusterIDWithQueryPayload) (*model.PaginatedResponse[link.Link], error) {
			userID := middleware.GetUserID(c)
			logger := middleware.GetLogger(c)
			return h.linkService.GetLinksByCusterID(c.Request().Context(), logger, userID, payload.ClusterID, payload.GetLinksQuery)
		},
		http.StatusOK,
		&link.GetLinkByClusterIDWithQueryPayload{},
	)(c)
}
