package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Adedunmol/glimpse/internal/errs"
	"github.com/Adedunmol/glimpse/internal/model"
	"github.com/Adedunmol/glimpse/internal/model/photo"
	"github.com/Adedunmol/glimpse/internal/model/upload"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UploadRepository struct {
	server *server.Server
}

func NewUploadRepository(server *server.Server) *UploadRepository {
	return &UploadRepository{
		server: server,
	}
}

func (r *UploadRepository) CreateUpload(ctx context.Context, userID string, payload *upload.CreateUploadPayload) (*upload.Upload, error) {
	stmt := `
		INSERT INTO uploads (name, expires_at, host_id)
		VALUES (@name, @expires_at, @host_id)
		RETURNING *
	`
	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"name":       payload.Name,
		"expires_at": payload.ExpiresAt,
		"host_id":    userID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute create upload query for user_id=%s name=%s: %w", userID, payload.Name, err)
	}

	uploadItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[upload.Upload])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:uploads for user_id=%s name=%s: %w", userID, payload.Name, err)
	}

	return &uploadItem, nil
}

func (r *UploadRepository) GetUploadByID(ctx context.Context, userID string, uploadID uuid.UUID) (*upload.PopulatedUpload, error) {

	stmt := `
			SELECT
				u.*,
				COALESCE(
					json_agg(p.* ORDER BY p.created_at) FILTER (WHERE p.id IS NOT NULL),
					'[]'
				) AS photos
			FROM uploads u
			LEFT JOIN photos p ON p.upload_id = u.id
			WHERE u.id = @id AND u.host_id = @host_id
			GROUP BY u.id
		`

	rows, err := r.server.DB.Pool.Query(ctx, stmt, pgx.NamedArgs{
		"id":      uploadID,
		"host_id": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get upload by id query for user_id=%s upload_id=%s: %w", userID, uploadID, err)
	}

	type populatedUploadRow struct {
		upload.Upload
		PhotosJSON json.RawMessage `db:"photos"`
	}

	row, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[populatedUploadRow])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:uploads for user_id=%s upload_id=%s: %w", userID, uploadID, err)
	}

	var photos []photo.Photo
	if err := json.Unmarshal(row.PhotosJSON, &photos); err != nil {
		return nil, fmt.Errorf("failed to unmarshal photos for upload_id=%s: %w", uploadID, err)
	}

	return &upload.PopulatedUpload{
		Upload: row.Upload,
		Photos: photos,
	}, nil
}

func (r *UploadRepository) GetUploads(ctx context.Context, userID string, query *upload.GetUploadsQuery) (*model.PaginatedResponse[upload.Upload], error) {

	stmt := `
		SELECT *
		FROM uploads u
	`
	args := pgx.NamedArgs{
		"host_id": userID,
	}
	conditions := []string{"host_id = @host_id"}

	if query.Status != nil {
		conditions = append(conditions, "status = @status")
		args["status"] = *query.Status
	}

	if query.Search != nil {
		conditions = append(conditions, "u.name ILIKE @search")
		args["search"] = "%" + *query.Search + "%"
	}

	if len(conditions) > 0 {
		stmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	countStmt := "SELECT COUNT(*) FROM uploads u"
	if len(conditions) > 0 {
		countStmt += " WHERE " + strings.Join(conditions, " AND ")
	}

	var total int
	err := r.server.DB.Pool.QueryRow(ctx, countStmt, args).Scan(&total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count for uploads user_id=%s: %w", userID, err)
	}

	if query.Sort != nil {
		stmt += " ORDER BY u." + *query.Sort
		if query.Order != nil && *query.Order == "desc" {
			stmt += " DESC "
		} else {
			stmt += " ASC "
		}
	} else {
		stmt += " ORDER BY u.created_at DESC "
	}
	stmt += " LIMIT @limit OFFSET @offset"
	args["limit"] = *query.Limit
	args["offset"] = (*query.Page - 1) * (*query.Limit)

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute get uploads query for user_id=%s: %w", userID, err)
	}

	todos, err := pgx.CollectRows(rows, pgx.RowToStructByName[upload.Upload])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return &model.PaginatedResponse[upload.Upload]{
				Data:       []upload.Upload{},
				Page:       *query.Page,
				Limit:      *query.Limit,
				Total:      0,
				TotalPages: 0,
			}, nil
		}
		return nil, fmt.Errorf("failed to collect rows from table:uploads for user_id=%s: %w", userID, err)
	}

	return &model.PaginatedResponse[upload.Upload]{
		Data:       todos,
		Page:       *query.Page,
		Limit:      *query.Limit,
		Total:      total,
		TotalPages: (total + *query.Limit - 1) / *query.Limit,
	}, nil
}

func (r *UploadRepository) UpdateUpload(ctx context.Context, userID string, payload *upload.UpdateUploadPayload) (*upload.Upload, error) {
	stmt := "UPDATE uploads SET "
	args := pgx.NamedArgs{
		"upload_id": payload.ID,
		"host_id":   userID,
	}
	setClauses := []string{}

	if payload.Name != nil {
		setClauses = append(setClauses, "name = @name")
		args["name"] = *payload.Name
	}

	if payload.ExpiresAt != nil {
		setClauses = append(setClauses, "expires_at = @expires_at")
		args["expires_at"] = *payload.ExpiresAt
	}

	if len(setClauses) == 0 {
		return nil, errs.NewBadRequestError("no fields to update", false, nil, nil, nil)
	}

	stmt += strings.Join(setClauses, ", ")
	stmt += " WHERE id = @upload_id AND host_id = @host_id RETURNING *"

	rows, err := r.server.DB.Pool.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update query for uploads for user_id=%s: %w", userID, err)
	}

	updatedUpload, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[upload.Upload])
	if err != nil {
		return nil, fmt.Errorf("failed to collect one row from table:uploads: %w", err)
	}

	return updatedUpload, nil
}

var UploadNotFound = "UPLOAD_NOT_FOUND"

func (r *UploadRepository) DeleteUpload(ctx context.Context, userID string, uploadID uuid.UUID) error {
	stmt := `
		DELETE FROM uploads
		WHERE
		id = @upload_id
		AND host_id = @user_id
	`
	result, err := r.server.DB.Pool.Exec(ctx, stmt, pgx.NamedArgs{
		"upload_id": uploadID,
		"user_id":   userID,
	})
	if err != nil {
		return fmt.Errorf("failed to execute delete query for uploads for user_id=%s: %w", userID, err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("upload not found", false, &UploadNotFound)
	}

	return nil
}
