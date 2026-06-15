package service

import (
	"context"

	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/cluster"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ClusterService struct {
	server      *server.Server
	clusterRepo *repository.ClusterRepository
}

func NewClusterService(srv *server.Server, clusterRepo *repository.ClusterRepository) *ClusterService {
	return &ClusterService{
		server:      srv,
		clusterRepo: clusterRepo,
	}
}

func (s *ClusterService) GetClusterByID(ctx context.Context, logger *zerolog.Logger, userID string, clusterID uuid.UUID) (*cluster.Cluster, error) {
	clusterItem, err := s.clusterRepo.GetClusterByID(ctx, userID, clusterID)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch cluster")
		return nil, err
	}

	return clusterItem, nil
}

func (s *ClusterService) GetClusters(ctx context.Context, logger *zerolog.Logger, userID string, query cluster.GetClustersQuery) (*model.PaginatedResponse[cluster.Cluster], error) {
	clusters, err := s.clusterRepo.GetClusters(ctx, userID, query)
	if err != nil {
		logger.Error().Err(err).Msg("failed to fetch clusters")
		return nil, err
	}

	return clusters, nil
}
