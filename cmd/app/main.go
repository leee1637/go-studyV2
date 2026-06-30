package main

import (
	"context"
	"fmt"
	"log"
	"os"
	repository_postgres "study/internal/features/users/repository/postgres"
	"study/internal/features/users/service"
	http_transport "study/internal/features/users/transport/http"
	"study/internal/features/users/transport/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("❌ Ошибка загрузки .env:", err)
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://test:123@127.0.0.1:5433/test_db?sslmode=disable"
	}
	fmt.Println("DEBUG DSN:", dsn)

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("❌ Ошибка создания пула подключений:", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatal("❌ БД не отвечает:", err)
	}
	fmt.Println("✅ Подключено к БД!")

	fmt.Println("hi")

	r := gin.Default()

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("❌ SECRET_KEY не найден в .env")
	}

	repo := repository_postgres.NewUserRepository(dbPool)
	service := service.NewAuthService(repo, secretKey)
	hand := http_transport.NewAuthHandler(service)

	r.POST("/register", hand.SignUp)
	r.POST("/login", hand.SignIn)

	api := r.Group("/api")

	api.Use(middleware.AuthMiddleware(secretKey))
	{

		api.GET("/admin/panel", middleware.RequireRole("ADMIN"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Добро пожаловать, Админ!"})
		})

		api.GET("/groups", middleware.RequireRole("ADMIN", "TEACHER"), func(c *gin.Context) {

		})
	}

	r.Run(":8081")

}
