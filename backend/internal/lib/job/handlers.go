package job

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Adedunmol/glimpse/internal/config"
	"github.com/Adedunmol/glimpse/internal/lib/email"
	"github.com/hibiken/asynq"
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

func (j *JobService) consumeMLStream(ctx context.Context) {

}
