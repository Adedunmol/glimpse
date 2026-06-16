package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/link"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LinkRepository struct {
	server *server.Server
}

func NewLinkRepository(srv *server.Server) *LinkRepository {
	return &LinkRepository{
		server: srv,
	}
}

func (l *LinkRepository) GetLinkByID(ctx context.Context, userID string, linkID uuid.UUID) (*link.Link, error) {
	stmt := `
		SELECT
			l.*
		FROM
			links l
		JOIN
			clusters c
		ON
			l.cluster_id = c.id
		JOIN
			uploads u
		ON
			c.upload_id = u.id
		WHERE
			l.id=@id AND u.host_id = @userID
	`

	args := pgx.NamedArgs{
		"id":     linkID,
		"userID": userID,
	}

	rows, err := l.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get link by id query for user_id=%s link_id=%s: %w", userID, linkID, err)
	}

	linkItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[link.Link])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:links for user_id=%s link_id=%s: %w", userID, linkID, err)
	}

	return &linkItem, nil
}

func (l *LinkRepository) GetLinkByToken(ctx context.Context, userID, token string) (*link.Link, error) {
	stmt := `
		SELECT
			l.*
		FROM
			links l
		JOIN
			clusters c
		ON
			l.cluster_id = c.id
		JOIN
			uploads u
		ON
			c.upload_id = u.id
		WHERE
			l.token=@token AND u.host_id = @userID
	`

	args := pgx.NamedArgs{
		"token":  token,
		"userID": userID,
	}

	rows, err := l.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get link by id query for user_id=%s token=%s: %w", userID, token, err)
	}

	linkItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[link.Link])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:links for user_id=%s token=%s: %w", userID, token, err)
	}

	return &linkItem, nil
}

func (l *LinkRepository) GetLinks(ctx context.Context, userID string, query *link.GetLinksQuery) (*model.PaginatedResponse[link.Link], error) {
	stmt := `
		SELECT
			l.*
		FROM
			links l
		JOIN
			clusters c
		ON
			l.cluster_id = c.id
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
		conditions = append(conditions, "l.token ILIKE @search")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	countStmt := `
		SELECT
			COUNT(*)
		FROM
			links l
		JOIN
			clusters c
		ON
			l.cluster_id = c.id
		JOIN
			uploads u
		ON
			c.upload_id = u.id
	`
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := l.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count for links user_id=%s: %w", userID, err)
	}

	if query.Sort != nil {
		stmt += " ORDER BY l." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC "
		} else {
			stmt += " ASC "
		}
	} else {
		stmt += " ORDER BY l.created_at DESC "
	}

	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := l.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get links query for user_id=%s: %w", userID, err)
	}

	linksData, err := pgx.CollectRows(rows, pgx.RowToStructByName[link.Link])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[link.Link]{
				Data:       []link.Link{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:links for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[link.Link]{
		Data:       linksData,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

func (l *LinkRepository) GetLinksByClusterID(ctx context.Context, userID string, clusterID uuid.UUID, query *link.GetLinksQuery) (*model.PaginatedResponse[link.Link], error) {
	stmt := `
		SELECT
			l.*
		FROM
			links l
		JOIN
			clusters c
		ON
			l.cluster_id = c.id
		JOIN
			uploads u
		ON
			c.upload_id = u.id
	`
	args := pgx.NamedArgs{
		"cluster_id": clusterID,
		"host_id":    userID,
	}

	conditions := []string{"u.host_id = @host_id", "c.cluster_id = @cluster_id"}

	if query.Search != nil {
		conditions = append(conditions, "l.token ILIKE @search")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	countStmt := `
		SELECT
			COUNT(*)
		FROM
			links l
		JOIN
			clusters c
		ON
			l.cluster_id = c.id
		JOIN
			uploads u
		ON
			c.upload_id = u.id
	`
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := l.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count for links user_id=%s: %w", userID, err)
	}

	if query.Sort != nil {
		stmt += " ORDER BY l." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC "
		} else {
			stmt += " ASC "
		}
	} else {
		stmt += " ORDER BY l.created_at DESC "
	}

	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := l.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get links query for user_id=%s: %w", userID, err)
	}

	linksData, err := pgx.CollectRows(rows, pgx.RowToStructByName[link.Link])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[link.Link]{
				Data:       []link.Link{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:links for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[link.Link]{
		Data:       linksData,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}
