package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/model/user"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/rs/zerolog"
)

const (
	UserCreatedEvent = "user.created"
	UserDeletedEvent = "user.deleted"
	UserUpdatedEvent = "user.updated"
)

type ClerkService struct {
	server   *server.Server
	userRepo repository.UserRepository
}

func NewClerkService(srv *server.Server, ur repository.UserRepository) *ClerkService {
	return &ClerkService{server: srv, userRepo: ur}
}

func (c *ClerkService) HandleNewUserEvent(ctx context.Context, logger *zerolog.Logger, payload user.CreateUserDTO) error {
	existing, err := c.userRepo.GetUserEmail(ctx, payload.Email)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("user.created: checking existing user: %w", err)
	}
	if existing != "" {
		logger.Warn().Str("event", "user_created").Msg("duplicate event, skipping...")
		return nil
	}

	if _, err = c.userRepo.CreateUser(ctx, payload.Email, payload.ClerkUserID); err != nil {
		logger.Error().Err(err).Msg("user.created: failed to create user")
		return fmt.Errorf("user.created: %w", err)
	}

	return nil
}

func (c *ClerkService) HandleDeleteUserEvent(ctx context.Context, userId string) error {
	return c.userRepo.DeleteUser(ctx, userId)
}
