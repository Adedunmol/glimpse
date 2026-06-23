package job

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Adedunmol/glimpse/internal/config"
	"github.com/Adedunmol/glimpse/internal/lib/email"
	"github.com/Adedunmol/glimpse/internal/model/upload"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

var emailClient *email.Client

func (j *JobService) InitHandlers(config *config.Config, logger *zerolog.Logger) {
	emailClient = email.NewClient(config, logger)
}

func (j *JobService) handleWelcomeEmailTask(ctx context.Context, t *asynq.Task) error {
	var p WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal welcome email payload: %w", err)
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Processing welcome email task")

	err := emailClient.SendWelcomeEmail(
		p.To,
		p.FirstName,
	)
	if err != nil {
		j.logger.Error().
			Str("type", "welcome").
			Str("to", p.To).
			Err(err).
			Msg("Failed to send welcome email")
		return err
	}

	j.logger.Info().
		Str("type", "welcome").
		Str("to", p.To).
		Msg("Successfully sent welcome email")
	return nil
}

func (j *JobService) handleLinkCleanup(ctx context.Context, t *asynq.Task) error {
	stmt := `UPDATE links SET is_active = false WHERE expires_at <= now() AND is_active = true`

	result, err := j.db.Exec(ctx, stmt)
	if err != nil {
		j.logger.Error().Err(err).Msg("failed to deactivate expired links")
		return fmt.Errorf("failed to deactivate expired links: %w", err)
	}

	j.logger.Info().
		Int64("deactivated", result.RowsAffected()).
		Msg("link cleanup completed")

	return nil
}

func (j *JobService) getUpload(ctx context.Context, uploadID string) (*upload.Upload, error) {
	stmt := `
		SELECT *
		FROM uploads
		WHERE id = @id
	`
	rows, err := j.db.Query(ctx, stmt, pgx.NamedArgs{
		"id": uploadID,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to execute get upload by id query upload_id=%s: %w", uploadID, err)
	}

	uploadItem, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[upload.Upload])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:uploads upload_id=%s: %w", uploadID, err)
	}

	return &uploadItem, nil
}

func (j *JobService) PublishNotification(ctx context.Context, taskType, userID, title, message string) error {
	task, err := NewNotificationTask(taskType, userID, title, message)
	if err != nil {
		j.logger.Error().Err(err).Msg("failed to create notifical task")
		return err
	}

	_, err = j.Client.EnqueueContext(ctx, task,
		asynq.Queue(CriticalPriority),
		asynq.MaxRetry(3),
	)

	return nil
}

func (j *JobService) handleNotificationTask(ctx context.Context, t *asynq.Task) error {
	var p NotificationPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("failed to unmarshal notification payload: %w", err)
	}

	j.logger.Info().
		Str("type", t.Type()).
		Str("userId", p.UserID).
		Msg("processing notification task")

	if err := j.notificationService.SendToUser(ctx, p.UserID, p.Title, p.Message); err != nil {
		j.logger.Error().
			Str("type", t.Type()).
			Str("userId", p.UserID).
			Err(err).
			Msg("failed to send notification")
		return err
	}

	j.logger.Info().
		Str("type", t.Type()).
		Str("userId", p.UserID).
		Msg("notification sent successfully")

	return nil
}

const (
	goConsumerGroup = "GoEventWorkerGroup"
)

func (j *JobService) consumeMLStream(ctx context.Context) {
	j.redisClient.XGroupCreateMkStream(ctx, j.streamName, goConsumerGroup, "0")

	for {
		select {
		case <-ctx.Done():
			return
		default:
			results, err := j.redisClient.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    goConsumerGroup,
				Consumer: "go-worker-1",
				Streams:  []string{j.streamName, ">"},
				Count:    10,
				Block:    5 * time.Second,
			}).Result()

			if err == redis.Nil {
				j.logger.Info().Msg("No new events")
				continue
			}

			if err != nil {
				j.logger.Error().Err(err).Msg("failed to read from ml stream")
				continue
			}

			for _, stream := range results {
				for _, msg := range stream.Messages {
					eventType, _ := msg.Values["type"].(string)

					if eventType == "cluster_event" {
						uploadID, _ := msg.Values["upload_id"].(string)
						j.logger.Info().Str("upload_id", uploadID).Msg("cluster event received")
						uploadItem, err := j.getUpload(ctx, uploadID)
						if err != nil {
							j.logger.Error().Err(err).Msgf("failed to fetch upload with upload_id=%s", uploadID)
							continue
						}

						if uploadItem == nil {
							j.logger.Error().Msgf("upload_id %s return null Upload item", uploadID)
							continue
						}

						if err := j.notificationService.SendToUser(ctx, uploadItem.HostID, "Your photos are ready", "Clustering is complete, check your album"); err != nil {
							j.logger.Error().Err(err).Str("upload_id", uploadID).Msg("failed to send cluster notification")
							continue
						}
					}

					j.redisClient.XAck(ctx, j.streamName, goConsumerGroup, msg.ID)
				}
			}
		}
	}
}
