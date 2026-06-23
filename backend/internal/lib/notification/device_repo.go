package notification

import (
	"context"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/errs"
	"github.com/Adedunmol/glimpse/internal/model/user_device"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DeviceNotFound = "DEVICE_NOT_FOUND"

type DeviceRepository struct {
	db *pgxpool.Pool
}

func NewDeviceRepository(db *pgxpool.Pool) *DeviceRepository {
	return &DeviceRepository{
		db: db,
	}
}

func (r *DeviceRepository) UpsertDevice(ctx context.Context, userID, fcmToken, platform string) (*user_device.UserDevice, error) {
	stmt := `
		INSERT INTO
			user_devices (user_id, push_token, platform)
		VALUES
			(@userId, @pushToken, @platform)
		ON CONFLICT
			(push_token)
		DO UPDATE SET
			user_id = @userId, platform = @platform
		RETURNING
			*
	`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"userId":    userID,
		"pushToken": fcmToken,
		"platform":  platform,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute create devuce query for user_id=%s: %w", userID, err)
	}

	deviceItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[user_device.UserDevice])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:device_token for user_id=%s: %w", userID, err)
	}
	return &deviceItem, nil
}

func (r *DeviceRepository) GetUserTokens(ctx context.Context, userID string) ([]user_device.UserDevice, error) {
	stmt := `
		SELECT
			*
		FROM
			user_devices
		WHERE
			user_id = @userId
		LIMIT 100
	`
	rows, err := r.db.Query(ctx, stmt, pgx.NamedArgs{
		"userId": userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute fetch user devices for user_id=%s", userID)
	}

	devices, err := pgx.CollectRows(rows, pgx.RowToStructByName[user_device.UserDevice])
	if err != nil {
		return nil, fmt.Errorf("failed to collect rows from table:user_devices for user_id=%s: %w", userID, err)
	}
	return devices, nil
}

func (r *DeviceRepository) DeleteDevice(ctx context.Context, userID, deviceToken string) error {
	stmt := `
		DELETE FROM
			user_devices
		WHERE
			user_id = @userId AND push_token = @deviceToken
	`

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"userId":      userID,
		"deviceToken": deviceToken,
	})

	if err != nil {
		return fmt.Errorf("failed to execute delete query for user device for user_id=%s: %w", userID, err)
	}

	if result.RowsAffected() == 0 {
		return errs.NewNotFoundError("device not found", false, &DeviceNotFound)
	}

	return nil
}

func (r *DeviceRepository) DeleteBulkDevice(ctx context.Context, userID string, deviceTokens []string) error {
	if len(deviceTokens) == 0 {
		return nil
	}

	stmt := `
        DELETE FROM
					user_devices
        WHERE
					user_id = @userId AND push_token = ANY(@deviceTokens)
    `

	result, err := r.db.Exec(ctx, stmt, pgx.NamedArgs{
		"userId":       userID,
		"deviceTokens": deviceTokens,
	})
	if err != nil {
		return fmt.Errorf("failed to execute bulk delete query for user_id=%s: %w", userID, err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("failed to delete device for user_id=%s", userID)
	}

	return nil
}
