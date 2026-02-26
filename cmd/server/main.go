package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"usermanagement/internal/application/user"
	"usermanagement/internal/delivery/http"
	"usermanagement/internal/infra/config"
	"usermanagement/internal/infra/logger"
	"usermanagement/internal/infra/persistence/postgres"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	// Initialize logger
	log, err := logger.New(cfg.Environment)
	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}
	defer log.Sync()

	log.Info("starting user management service",
		zap.String("environment", cfg.Environment),
		zap.String("port", cfg.HTTPPort),
	)

	// Database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer pool.Close()

	// Test database connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatal("failed to ping database", zap.Error(err))
	}
	log.Info("connected to database")

	// Dependency Injection
	// Infra
	userRepo := postgres.NewUserRepository(pool, log)

	// Application (Use Cases)
	createUC := user.NewCreateUserUseCase(userRepo)
	getUC := user.NewGetUserUseCase(userRepo)
	listUC := user.NewListUsersUseCase(userRepo)
	updateUC := user.NewUpdateUserUseCase(userRepo)
	deleteUC := user.NewDeleteUserUseCase(userRepo)

	// Delivery
	handler := http.NewUserHandler(createUC, getUC, listUC, updateUC, deleteUC, log)
	router := http.NewRouter(handler, log)

	// HTTP Server
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info("server is shutting down...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			log.Fatal("could not gracefully shutdown the server", zap.Error(err))
		}
		close(done)
	}()

	log.Info("server is ready to handle requests", zap.String("addr", srv.Addr))

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("could not listen on", zap.String("addr", srv.Addr), zap.Error(err))
	}

	<-done
	log.Info("server stopped")
}