package service

import (
	"context"
	"time"

	custom_bcrypt "github.com/Adedunmol/glimpse/internal/lib/bcrypt"
	"github.com/Adedunmol/glimpse/internal/lib/utils"
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/cluster"
	"github.com/Adedunmol/glimpse/internal/model/link"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type LinkService struct {
	server   *server.Server
	linkRepo *repository.LinkRepository
}

func NewLinkService(srv *server.Server, linkRepo *repository.LinkRepository) *LinkService {
	return &LinkService{
		server:   srv,
		linkRepo: linkRepo,
	}
}

func (s *LinkService) CreateLink(ctx context.Context, logger *zerolog.Logger, userID string, payload *cluster.CreateLinkCommand) (*link.Link, error) {
	linkToken := utils.GenerateToken()

	var passwordHash *string
	isPasswordProtected := false
	if payload.Password != "" {
		hash, err := custom_bcrypt.HashPassword(payload.Password)
		if err != nil {
			logger.Error().Err(err).Msg("failed to hash password")
			return nil, err
		}
		passwordHash = &hash
		isPasswordProtected = true
	}

	expiresAt := time.Now().Add(time.Hour * 24 * 30) // 30 days
	if payload.ExpiresAt != nil {
		expiresAt = *payload.ExpiresAt
	}

	isActive := true
	linkPayload := &link.CreateLinkPayload{
		ClusterID:           payload.ClusterID,
		Token:               linkToken,
		IsPasswordProtected: &isPasswordProtected,
		PasswordHash:        passwordHash,
		ExpiresAt:           &expiresAt,
		IsActive:            &isActive,
	}

	linkItem, err := s.linkRepo.CreateLink(ctx, userID, linkPayload)
	if err != nil {
		logger.Error().Err(err).Msg("failed to create link item")
		return nil, err
	}

	return linkItem, nil
}

func (s *LinkService) GetLinkByID(ctx context.Context, logger *zerolog.Logger, userID string, linkID uuid.UUID) (*link.Link, error) {
	linkItem, err := s.linkRepo.GetLinkByID(ctx, userID, linkID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch link")
		return nil, err
	}

	return linkItem, nil
}

func (s *LinkService) GetLinkByToken(ctx context.Context, logger *zerolog.Logger, userID, token string) (*link.Link, error) {
	linkItem, err := s.linkRepo.GetLinkByToken(ctx, userID, token)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch link")
		return nil, err
	}

	return linkItem, nil
}

func (s *LinkService) GetLinks(ctx context.Context, logger *zerolog.Logger, userID string, query *link.GetLinksQuery) (*model.PaginatedResponse[link.Link], error) {
	links, err := s.linkRepo.GetLinks(ctx, userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch links")
		return nil, err
	}

	return links, err
}

func (s *LinkService) GetLinksByCusterID(ctx context.Context, logger *zerolog.Logger, userID string, clusterID uuid.UUID, query *link.GetLinksQuery) (*model.PaginatedResponse[link.Link], error) {
	links, err := s.linkRepo.GetLinksByClusterID(ctx, userID, clusterID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch links")
		return nil, err
	}

	return links, err
}
