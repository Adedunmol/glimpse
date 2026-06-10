package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
)

const (
	userCreatedEvent = "user.created"
	userDeletedEvent = "user.deleted"
	userUpdatedEvent = "user.updated"
)

type ClerkService struct {
	server   *server.Server
	userRepo repository.UserRepository
}

type ClerkEventPayload struct {
	Data       json.RawMessage `json:"data"`
	Object     string          `json:"object"`
	EventType  string          `json:"type"`
	TimeStamp  string          `json:"timestamp"`
	InstanceID string          `json:"instance_id"`
}

// TODO: implement validator logic for payload
func (p *ClerkEventPayload) Validate() error {
	return nil
}

type clerkUserData struct {
	ID             string              `json:"id"`
	EmailAddresses []clerkEmailAddress `json:"email_addresses"`
}
type clerkEmailAddress struct {
	EmailAddress string `json:"email_address"`
}
type clerkDeletedData struct {
	ID      string `json:"id"`
	Deleted bool   `json:"deleted"`
}

func NewClerkService(srv *server.Server, ur repository.UserRepository) *ClerkService {
	return &ClerkService{server: srv, userRepo: ur}
}

func (c *ClerkService) HandleClerkEvents(ctx context.Context, payload ClerkEventPayload) error {

	switch payload.EventType {
	case userCreatedEvent:
		c.server.Logger.Info().Msg("clerk event: user.created")

		var userData clerkUserData
		if err := json.Unmarshal(payload.Data, &userData); err != nil {
			return fmt.Errorf("user.created: malformed data: %w", err)
		}
		if len(userData.EmailAddresses) == 0 {
			return fmt.Errorf("user.created: no email addresses in payload")
		}

		email := userData.EmailAddresses[0].EmailAddress

		existing, err := c.userRepo.GetUserEmail(ctx, email)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return fmt.Errorf("user.created: checking existing user: %w", err)
		}
		if existing != "" {
			c.server.Logger.Warn().Msg("user.created: duplicate event, skipping")
			return nil
		}

		if _, err = c.userRepo.CreateUser(ctx, email, userData.ID); err != nil {
			c.server.Logger.Error().Err(err).Msg("user.created: failed to create user")
			return fmt.Errorf("user.created: %w", err)
		}
	case userDeletedEvent:
		c.server.Logger.Info().Msg("clerk event: user.deleted")

		var userData clerkDeletedData
		if err := json.Unmarshal(payload.Data, &userData); err != nil {
			return fmt.Errorf("user.deleted: malformed data: %w", err)
		}

		if err := c.userRepo.DeleteUser(ctx, userData.ID); err != nil {
			c.server.Logger.Error().Err(err).Msg("user.deleted: failed to delete user")
			return fmt.Errorf("user.deleted: %w", err)
		}
	case userUpdatedEvent:
		c.server.Logger.Info().Msg("clerk event: user.updated (unimplemented)")

	default:
		c.server.Logger.Warn().Str("event_type", payload.EventType).Msg("unhandled clerk event")
	}
	return nil
}
