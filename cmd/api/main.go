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

	shutdownCh := make(chan struct{})

	go func() {
		sigintCh := make(chan os.Signal, 1)
		signal.Notify(sigintCh, os.Interrupt, syscall.SIGTERM)

		slog.Info("Server shutting down", "addr", server.Addr)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.TimeoutIdle)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			slog.Error("Shutting down server failure", "msg", err)
		}

		db, err := db.DB()
		if err == nil {
			if err = db.Close(); err != nil {
				slog.Error("failed to close db connection")
			}
		}

		close(shutdownCh)
	}()

	slog.Info("Starting server ...", "ADDRESS", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup failed")
	}

	<-shutdownCh

	slog.Info("server shutdown success")
}
