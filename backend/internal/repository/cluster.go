package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/cluster"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ClusterRepository struct {
	server *server.Server
}

func NewClusterRepository(srv *server.Server) *ClusterRepository {
	return &ClusterRepository{
		server: srv,
	}
}

func (c *ClusterRepository) GetClusterByID(ctx context.Context, userID string, clusterID uuid.UUID) (*cluster.Cluster, error) {
	stmt := `
		SELECT
			c.*
		FROM
			clusters c
		JOIN
			uploads u
		ON
			c.upload_id = u.id
		WHERE
			c.id = @id AND u.host_id = @userID
	`
	args := pgx.NamedArgs{
		"id":     clusterID,
		"userID": userID,
	}

	rows, err := c.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get cluster by id query for user_id=%s cluster_id=%s: %w", userID, clusterID, err)
	}

	clusterItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[cluster.Cluster])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:clusters for user_id=%s cluster_id=%s: %w", userID, clusterID, err)
	}

	return &clusterItem, nil
}

func (c *ClusterRepository) GetClusters(ctx context.Context, userID string, query *cluster.GetClustersQuery) (*model.PaginatedResponse[cluster.Cluster], error) {
	stmt := `
		SELECT
			c.*
		FROM
			clusters c
		JOIN
			uploads u
		ON
			c.upload_id = u.id
	`
	args := pgx.NamedArgs{
		"host_id": userID,
	}
	conditions := []string{"u.host_id = @host_id"}

	if query.Search != nil {
		conditions = append(conditions, "c.label ILIKE @search")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	countStmt := `
		SELECT
			COUNT(*)
		FROM
			clusters c
		JOIN
			uploads u
		ON
			c.upload_id = u.id
	`
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := c.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count for clusters user_id=%s: %w", userID, err)
	}

	if query.Sort != nil {
		stmt += " ORDER BY c." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC "
		} else {
			stmt += " ASC "
		}
	} else {
		stmt += " ORDER BY c.created_at DESC "
	}

	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := c.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get clusters query for user_id=%s: %w", userID, err)
	}

	clusters, err := pgx.CollectRows(rows, pgx.RowToStructByName[cluster.Cluster])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[cluster.Cluster]{
				Data:       []cluster.Cluster{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:clusters for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[cluster.Cluster]{
		Data:       clusters,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}
