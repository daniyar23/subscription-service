package main

import (
	"fmt"
	"log"

	"github.com/daniyar23/subscribe-service/internal/handler"
	"github.com/daniyar23/subscribe-service/internal/repository/postgres"
	"github.com/daniyar23/subscribe-service/internal/service"
	"github.com/daniyar23/subscribe-service/pkg/database"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Подгружаем локальные переменные окружения (чтобы не брались системные)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	//  СОЗДАЕМ ПУЛЛ ПОДЛКЮЧЕНИЙ К БД
	if err := database.ConnectDB(); err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// ОТКЛАДВАЕМ ОТКЛЮЧЕНИЕ СОЕДИНЕНИЕ К БД
	defer database.DisconnectDB()
	fmt.Println("Database connected successfully!")

	// СОЗДАЕМ НАШИ СЛОИ
	repo := postgres.NewSubscriptionRepo(database.Pool)
	service := service.NewSubscriptionService(repo)
	h := handler.NewSubscriptionHandler(service)

	router := gin.Default()
	handler.RegisterRoutes(router, h)
	router.Run(":8080")

}
