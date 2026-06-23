package service

import (
	"context"

	"github.com/Adedunmol/glimpse/internal/errs"
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/cluster"
	"github.com/Adedunmol/glimpse/internal/model/link"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ClusterService struct {
	server      *server.Server
	clusterRepo *repository.ClusterRepository
	linkService *LinkService
}

func NewClusterService(srv *server.Server, clusterRepo *repository.ClusterRepository, link *LinkService) *ClusterService {
	return &ClusterService{
		server:      srv,
		clusterRepo: clusterRepo,
		linkService: link,
	}
}

func (s *ClusterService) CreateClusterLink(ctx context.Context, logger *zerolog.Logger, userID string, payload *cluster.CreateLinkCommand) (*link.Link, error) {
	clusterItem, err := s.GetClusterByID(ctx, logger, userID, payload.ClusterID)
	if err != nil {
		logger.Error().Err(err).Msgf("failed to fetch cluster cluster_id=%s user_id=%s", payload.ClusterID.String(), userID)
		return nil, err
	}
	if clusterItem == nil {
		logger.Warn().Msgf("cluster not found cluster_id=%s user_id=%s", payload.ClusterID.String(), userID)
		return nil, errs.NewNotFoundError("cluster item not found", false, nil)
	}

	return s.linkService.CreateLink(ctx, logger, userID, payload)
}

func (s *ClusterService) GetClusterByID(ctx context.Context, logger *zerolog.Logger, userID string, clusterID uuid.UUID) (*cluster.Cluster, error) {
	clusterItem, err := s.clusterRepo.GetClusterByID(ctx, userID, clusterID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch cluster")
		return nil, err
	}

	return clusterItem, nil
}

func (s *ClusterService) GetClusters(ctx context.Context, logger *zerolog.Logger, userID string, query *cluster.GetClustersQuery) (*model.PaginatedResponse[cluster.Cluster], error) {
	clusters, err := s.clusterRepo.GetClusters(ctx, userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch clusters")
		return nil, err
	}

	return clusters, nil
}
