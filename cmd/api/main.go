package main

import (
	"fmt"
	"log"

	order_usecase "github.com/Disneyjr/OrderAPI-CleanArch/internal/application"
	configs "github.com/Disneyjr/OrderAPI-CleanArch/internal/config"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/domain"
	orderRepository "github.com/Disneyjr/OrderAPI-CleanArch/internal/infrastructure/persistence/postgres"
	"github.com/Disneyjr/OrderAPI-CleanArch/internal/infrastructure/rest"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("Starting APP!")
	defer fmt.Println("Stoped APP!")
	configs, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configs.DBHost, configs.DBPort, configs.DBUser, configs.DBPassword, configs.DBName,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&domain.Order{})
	repo := orderRepository.NewOrderRepository(db)
	useCase := order_usecase.NewOrderUseCase(repo)
	orderHandler := rest.NewOrderHandler(useCase)
	r := gin.Default()
	r.GET("/health", health)
	r.GET("/orders", orderHandler.GetOrders)
	r.GET("/orders/:id", orderHandler.GetOrder)
	r.POST("/orders", orderHandler.Create)
	r.Run("localhost:" + configs.WebServerPort)
}

func health(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}
