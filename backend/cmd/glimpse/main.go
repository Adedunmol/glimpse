package main

import (
	"context"
	"errors"
	"net/http"
	_ "net/http/pprof" // side-effect: registers /debug/pprof routes
	"os"
	"os/signal"

	"time"

	"github.com/Adedunmol/glimpse/internal/config"
	"github.com/Adedunmol/glimpse/internal/database"
	"github.com/Adedunmol/glimpse/internal/handler"
	"github.com/Adedunmol/glimpse/internal/lib/fcm"
	"github.com/Adedunmol/glimpse/internal/lib/job"
	"github.com/Adedunmol/glimpse/internal/lib/notification"
	"github.com/Adedunmol/glimpse/internal/logger"
	"github.com/Adedunmol/glimpse/internal/repository"
	"github.com/Adedunmol/glimpse/internal/router"
	"github.com/Adedunmol/glimpse/internal/server"
	"github.com/Adedunmol/glimpse/internal/service"
)

const DefaultContextTimeout = 30

func main() {

	// go func() {
	// 	l.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	// f, _ := os.Create("cpu.prof")
	// defer f.Close()

	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()

	cfg, err := config.LoadConfig()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize New Relic logger service
	loggerService := logger.NewLoggerService(cfg.Observability)
	defer loggerService.Shutdown()

	log := logger.NewLoggerWithService(cfg.Observability, loggerService)

	if cfg.Primary.Env != "local" {
		if err := database.Migrate(context.Background(), &log, cfg); err != nil {
			log.Fatal().Err(err).Msg("failed to migrate database")
		}
	}

	// Initialize server
	srv, err := server.New(cfg, &log, loggerService)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize server")
	}

	// Initialize repositories, services, and handlers
	fcmClient, err := fcm.NewFCMClient(cfg.FCM.CredentialPath, cfg.FCM.ProjectID)
	if err != nil {
		log.Warn().Msgf("failed to initialize fcm client: %s", err)
		fcmClient = nil
		//continue startup
	}
	deviceRepo := notification.NewDeviceRepository(srv.DB.Pool)
	repos := repository.NewRepositories(srv, deviceRepo)

	deviceService := notification.NewDeviceService(deviceRepo)
	notificationService := notification.NewNotificationService(deviceRepo, fcmClient)
	jobService := job.NewJobService(srv.Logger, cfg, srv.DB.Pool, srv.Redis, notificationService, cfg.Redis.StreamName)

	services, serviceErr := service.NewServices(srv, repos, deviceService, jobService)

	if serviceErr != nil {
		log.Fatal().Err(serviceErr).Msg("could not create services")
	}
	handlers := handler.NewHandlers(srv, services, deviceService)

	// Initialize router
	r := router.NewRouter(srv, handlers, services)

	// Setup HTTP server
	srv.SetupHTTPServer(r)

	//initialize asynq job service
	jobService.InitHandlers(cfg, srv.Logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Start server
	go func() {
		if err = srv.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("failed to start server")
		}
	}()

	// Start job service
	go func() {
		if err = jobService.Start(); err != nil {
			log.Fatal().Err(err).Msg("failed to start job service")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), DefaultContextTimeout*time.Second)
	defer cancel()

	// stop job service
	defer jobService.Stop()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("server forced to shutdown")
	}

	log.Info().Msg("server exited properly")
}
