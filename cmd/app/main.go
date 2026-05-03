package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/daniyar23/subscribe-service/internal/logger"
	"github.com/daniyar23/subscribe-service/internal/middleware"
	"go.uber.org/zap"

	"github.com/daniyar23/subscribe-service/internal/handler"
	"github.com/daniyar23/subscribe-service/internal/repository/postgres"
	"github.com/daniyar23/subscribe-service/internal/service"
	"github.com/daniyar23/subscribe-service/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// logger zap
	logg, err := logger.New()
	if err != nil {
		panic(err)
	}
	defer logg.Sync()

	// env
	if err := godotenv.Load(); err != nil {
		logg.Warn("no .env file found")
	}

	// DB connect
	if err := database.ConnectDB(); err != nil {
		logg.Fatal("failed to connect DB", zap.Error(err))
	}
	defer database.DisconnectDB()

	// layers
	repo := postgres.NewSubscriptionRepo(database.Pool, logg)
	service := service.NewSubscriptionService(repo, logg)
	h := handler.NewSubscriptionHandler(service, logg)

	// router
	router := gin.New() //(чтобы самому управлять middleware)
	router.Use(gin.Recovery())
	router.Use(middleware.Logger(logg))
	handler.RegisterRoutes(router, h)

	// server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// server start
	go func() {
		logg.Info("server started", zap.String("addr", ":8080"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logg.Fatal("listen failed", zap.Error(err))
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// graceful shutdown
	if err := srv.Shutdown(ctx); err != nil {
		logg.Error("server forced to shutdown", zap.Error(err))
	}

	logg.Info("server exited properly")
}
