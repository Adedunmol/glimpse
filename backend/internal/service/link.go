package service

import (
	"context"

	"github.com/Adedunmol/glimpse/internal/model"
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

func NewLinkRepository(srv *server.Server, linkRepo *repository.LinkRepository) *LinkService {
	return &LinkService{
		server:   srv,
		linkRepo: linkRepo,
	}
}

func (s *LinkService) GetLinkByID(ctx context.Context, logger *zerolog.Logger, userID string, linkID uuid.UUID) (*link.Link, error) {
	linkItem, err := s.linkRepo.GetLinkByID(ctx, userID, linkID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch link")
		return nil, err
	}

	return linkItem, nil
}

func (s *LinkService) GetLinks(ctx context.Context, logger *zerolog.Logger, userID string, query link.GetLinksQuery) (*model.PaginatedResponse[link.Link], error) {
	links, err := s.linkRepo.GetLinks(ctx, userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch links")
		return nil, err
	}

	return links, err
}
