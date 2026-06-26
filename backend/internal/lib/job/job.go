package job

import (
	"context"
	"fmt"
	"time"

	"github.com/Adedunmol/glimpse/internal/config"
	"github.com/Adedunmol/glimpse/internal/lib/notification"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

const (
	TaskLinkCleanup  = "link:cleanup"
	CriticalPriority = "critical"
	DefaultPriority  = "default"
	LowPriority      = "low"
)

type JobService struct {
	Client              *asynq.Client
	server              *asynq.Server
	scheduler           *asynq.Scheduler
	logger              *zerolog.Logger
	db                  *pgxpool.Pool
	notificationService *notification.NotificationService
	redisClient         *redis.Client
	streamName          string
	ctx                 context.Context
	cancelFunc          context.CancelFunc
}

func NewJobService(logger *zerolog.Logger, cfg *config.Config, pool *pgxpool.Pool, redisClient *redis.Client, notification *notification.NotificationService, streamName string) *JobService {
	redisAddr := cfg.Redis.Address
	redisOpts := asynq.RedisClientOpt{
		Addr: redisAddr,
	}

	client := asynq.NewClient(redisOpts)

	server := asynq.NewServer(
		redisOpts,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				CriticalPriority: 6, // Higher priority queue for important emails
				DefaultPriority:  3, // Default priority queue for most emails
				LowPriority:      1, // Lower priority queue for non-urgent emails
			},
		},
	)

	scheduler := asynq.NewScheduler(redisOpts, nil)

	ctx, cancelFunc := context.WithCancel(context.Background())

	return &JobService{
		Client:              client,
		server:              server,
		scheduler:           scheduler,
		logger:              logger,
		db:                  pool,
		notificationService: notification,
		redisClient:         redisClient,
		streamName:          streamName,
		ctx:                 ctx,
		cancelFunc:          cancelFunc,
	}
}

func (j *JobService) Start() error {
	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, j.handleWelcomeEmailTask)
	mux.HandleFunc(TaskLinkCleanup, j.handleLinkCleanup)
	mux.HandleFunc(TaskNotifyLinkCreated, j.handleNotificationTask)
	mux.HandleFunc(TaskNotifyUploadDone, j.handleNotificationTask)
	mux.HandleFunc(TaskNotifyClusterReady, j.handleNotificationTask)

	if _, err := j.scheduler.Register("@daily", asynq.NewTask(TaskLinkCleanup, nil)); err != nil {
		return fmt.Errorf("failed to register link cleanup job: %w", err)
	}

	j.logger.Info().Msg("Starting background job server")
	if err := j.server.Start(mux); err != nil {
		return err
	}

	if err := j.scheduler.Start(); err != nil {
		return fmt.Errorf("failed to start scheduler: %w", err)
	}

	go j.consumeMLStream(j.ctx)

	return nil
}

func (j *JobService) Stop() {
	j.logger.Info().Msg("Stopping background job server")

	if j.cancelFunc != nil {
		j.cancelFunc()
	}

	j.scheduler.Shutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	shutdownChan := make(chan struct{})
	go func() {
		j.server.Shutdown()
		close(shutdownChan)
	}()

	select {
	case <-shutdownChan:
		j.logger.Info().Msg("Job server stopped gracefully")
	case <-ctx.Done():
		j.logger.Warn().Msg("Job server shutdown timed out, forcing stop")
		j.server.Stop()
	}

	if err := j.Client.Close(); err != nil {
		j.logger.Error().Err(err).Msg("Failed to close Asynq client")
	}
}
