package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"study/internal/core/middleware"
	service_email "study/internal/features/email/service"
	students_postgres_repository "study/internal/features/students/repository/postgres"
	student_service "study/internal/features/students/service"
	http_student "study/internal/features/students/transport/http"
	repository_teacher "study/internal/features/teachers/repository"
	service_teacher "study/internal/features/teachers/service"
	transport_teacher "study/internal/features/teachers/transport"
	repository_postgres "study/internal/features/users/repository/postgres"
	"study/internal/features/users/service"
	http_transport "study/internal/features/users/transport/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Ошибка загрузки .env:", err)
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		dsn = "postgres://test:123@127.0.0.1:5534/test_db?sslmode=disable"
	}
	fmt.Println("DEBUG DSN:", dsn)

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("Ошибка создания пула подключений:", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		log.Fatal("БД не отвечает:", err)
	}
	fmt.Println("Подключено к БД!")

	fmt.Println("hi")

	r := gin.Default()

	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		log.Fatal("SECRET_KEY не найден в .env")
	}

	emailConfig := service_email.EmailConfig{
		Host:     "localhost",
		Port:     1025,
		Username: "",
		Password: "",
		From:     "sgty@university.ru",
		BaseURL:  "http://localhost:8091",
	}
	see := service_email.NewEmailConfig(emailConfig)
	repo := repository_postgres.NewUserRepository(dbPool)
	service := service.NewAuthService(repo, secretKey, *see)
	hand := http_transport.NewAuthHandler(service)

	studentRepo := students_postgres_repository.NewUserRepository(dbPool)
	studentService := student_service.NewStudentService(studentRepo, secretKey)
	studentHand := http_student.NewStudentHandler(studentService)

	teacherRepo := repository_teacher.NewTeacherRepository(dbPool)
	teacherService := service_teacher.NewTeacherService(teacherRepo, secretKey)
	teacherHand := transport_teacher.NewTeacherHandler(teacherService)

	r.POST("/api/auth/register", hand.SignUp)
	r.POST("/api/auth/login", hand.SignIn)
	r.GET("/api/auth/verify/:token", hand.VerifyEmail)
	r.POST("/api/auth/refresh", hand.Refresh)
	r.POST("/api/auth/logout", hand.Logout)

	api := r.Group("/api")

	api.Use(middleware.AuthMiddleware(secretKey))
	{

		api.GET("/admin/panel", middleware.RequireRole("ADMIN"), func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Добро пожаловать, Админ!"})
		})

		api.GET("/groups", middleware.RequireRole("ADMIN", "TEACHER"), func(c *gin.Context) {

		})

		students := api.Group("/students")
		{
			students.GET("", middleware.RequireRole("ADMIN", "TEACHER", "STUDENT"), studentHand.GetAll)
			students.GET("/:id", middleware.RequireRole("ADMIN", "TEACHER", "STUDENT"), studentHand.GetByID)
			students.GET("/group/:group", middleware.RequireRole("ADMIN", "TEACHER", "STUDENT"), studentHand.GetByGroup)
			students.DELETE("/:id", middleware.RequireRole("ADMIN", "TEACHER", "STUDENT"), studentHand.DeleteStudent)
			students.PATCH("/:id", middleware.RequireRole("ADMIN", "TEACHER", "STUDENT"), studentHand.UpdateStudent)

		}
		teachers := api.Group("/teachers")
		{
			teachers.GET("", middleware.RequireRole("ADMIN"), teacherHand.GetAll)
			teachers.GET("/:id", middleware.RequireRole("ADMIN", "TEACHER"), teacherHand.GetByID)
			teachers.PATCH("/:id", middleware.RequireRole("ADMIN", "TEACHER"), teacherHand.UpdateTeacher)
			teachers.DELETE("/:id", middleware.RequireRole("ADMIN"), teacherHand.DeleteTeacher)
			teachers.GET("/group/:group", middleware.RequireRole("ADMIN", "TEACHER"), teacherHand.GetByGroup)
			teachers.POST("/:id/group/:group", middleware.RequireRole("ADMIN"), teacherHand.AddGroup)
			teachers.DELETE("/:id/group/:group", middleware.RequireRole("ADMIN"), teacherHand.DeleteGroup)
		}

	}
	r.Run(":8091")
}
