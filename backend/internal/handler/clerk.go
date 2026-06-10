package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
	"github.com/labstack/echo/v4"
	svix "github.com/svix/svix-webhooks/go"
)

type ClerkWebHookHandler struct {
	Handler
	clerkService *service.ClerkService
	wh           *svix.Webhook
}

func NewClerkWebHookHandler(s *server.Server, cs *service.ClerkService) (*ClerkWebHookHandler, error) {
	wh, err := svix.NewWebhook(s.Config.Clerk.WebHookAuthorizationSecret)
	if err != nil {
		return nil, err
	}
	return &ClerkWebHookHandler{
		Handler:      NewHandler(s),
		clerkService: cs,
		wh:           wh,
	}, nil
}

func (h *ClerkWebHookHandler) HandleEvent(c echo.Context) error {
	return Handle(
		h.Handler,

		func(c echo.Context, _ *service.ClerkEventPayload) (any, error) {

			raw, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "failed to read body")
			}

			c.Request().Body = io.NopCloser(bytes.NewBuffer(raw))

			headers := c.Request().Header
			if err = h.wh.Verify(raw, headers); err != nil {
				h.server.Logger.Error().Err(err).Msg("failed to verify webhook signature")
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook signature")
			}

			var eventPayload service.ClerkEventPayload
			if err = json.Unmarshal(raw, &eventPayload); err != nil {
				h.server.Logger.Error().Err(err).Msg("error decoding raw payload")
				return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, "malformed json")
			}

			err = h.clerkService.HandleClerkEvents(c.Request().Context(), eventPayload)
			if err != nil {
				return nil, err
			}

			return nil, nil
		},
		http.StatusOK,
		&service.ClerkEventPayload{},
	)(c)
}
