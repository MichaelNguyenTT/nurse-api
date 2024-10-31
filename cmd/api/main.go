package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"nms/api/router"
	"nms/config"
	"os"
	"os/signal"
	"syscall"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const dbConn = "host=%s user=%s password=%s dbname=%s port=%d sslmode=disable"

func main() {
	cfg := config.LoadConfig()

	fmtDB := fmt.Sprintf(dbConn, cfg.DB.Host, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(fmtDB))
	if err != nil {
		log.Fatal("failed database connection")
		return
	}

	mux := router.NewRouter(db)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      mux,
		ReadTimeout:  cfg.Server.TimeoutRead,
		WriteTimeout: cfg.Server.TimeoutWrite,
		IdleTimeout:  cfg.Server.TimeoutIdle,
	}

	closeChan := make(chan os.Signal, 1)
	signal.Notify(closeChan, os.Interrupt, syscall.SIGTERM)
	serverErrors := make(chan error, 1)

	go func() {
		slog.Info("Starting server...", "addr", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case serverError := <-serverErrors:
		if serverError != nil && serverError != http.ErrServerClosed {
			slog.Error("Server error encounter", "error", err)
		}
	case <-closeChan:
		slog.Info("Server shutting down", "addr", server.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.TimeoutIdle)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			slog.ErrorContext(ctx, "Shut down server failure", "msg", err)

			if err := server.Close(); err != nil {
				slog.ErrorContext(ctx, "Couldn't close server", "msg", err)
			}
		}

		db, err := db.DB()
		if err == nil {
			if err = db.Close(); err != nil {
				slog.Error("failed to close db connection")
			}
		}
	}

	slog.Info("server shutdown success")
}
