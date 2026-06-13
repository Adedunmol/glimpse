package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Adedunmol/glimpse/internal/middleware"
	"github.com/Adedunmol/glimpse/internal/model/user"
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

		func(c echo.Context, _ *user.ClerkEventPayload) (any, error) {
			logger := middleware.GetLogger(c)
			raw, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return nil, echo.NewHTTPError(http.StatusBadRequest, "failed to read body")
			}

			c.Request().Body = io.NopCloser(bytes.NewBuffer(raw))

			headers := c.Request().Header
			if err = h.wh.Verify(raw, headers); err != nil {
				logger.Error().Err(err).Msg("failed to verify webhook signature")
				return nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid webhook signature")
			}

			var payload user.ClerkEventPayload
			if err = json.Unmarshal(raw, &payload); err != nil {
				logger.Error().Err(err).Msg("error decoding raw payload")
				return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, "malformed json")
			}

			logger.Info().
				Str("event", "clerk event").
				Str("event_type", payload.EventType).
				Msg("Clerk event received")

			switch payload.EventType {
			case service.UserCreatedEvent:
				var userData user.ClerkUserData
				if err := json.Unmarshal(payload.Data, &userData); err != nil {
					return nil, fmt.Errorf("user.created: malformed data: %w", err)
				}
				if len(userData.EmailAddresses) == 0 {
					return nil, fmt.Errorf("user.created: no email addresses in payload")
				}

				userDTO := user.CreateUserDTO{
					Email:       userData.EmailAddresses[0].EmailAddress,
					ClerkUserID: userData.ID,
				}

				if err := h.clerkService.HandleNewUserEvent(c.Request().Context(), logger, userDTO); err != nil {
					return nil, fmt.Errorf("failed to handle create event: %w", err)
				}
			case service.UserDeletedEvent:
				var userData user.ClerkDeletedData
				if err := json.Unmarshal(payload.Data, &userData); err != nil {
					return nil, fmt.Errorf("user.deleted: malformed data: %w", err)
				}

				if err := h.clerkService.HandleDeleteUserEvent(c.Request().Context(), userData.ID); err != nil {
					return nil, fmt.Errorf("failed to handle delete event: %w", err)
				}
			default:
				logger.Warn().Msg("unhandled clerk event")
			}
			return nil, nil
		},
		http.StatusOK,
		&user.ClerkEventPayload{},
	)(c)
}
