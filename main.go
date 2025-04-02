package main

import (
	"log"
	"mods/command"
	"mods/config"
	"mods/controller"
	"mods/middleware"
	"mods/repository"
	"mods/routes"
	"mods/service"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)
func args(db *gorm.DB) bool {
	if len(os.Args) > 1 {
		flag := command.Commands(db)
		if !flag {
			return false
		}
	}

	return true
}

func main() {
	db := config.SetUpDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	if !args(db) {
		return
	}

	var (
		jwtService service.JWTService = service.NewJWTService()

		// Implementation Dependency Injection
		// Repository
		userRepository repository.UserRepository = repository.NewUserRepository(db)
		examRepository repository.ExamRepository = repository.NewExamRepository(db)
		userExamRepository repository.UserExamRepository = repository.NewUserExamRepository(db)

		// Service
		userService service.UserService = service.NewUserService(userRepository, jwtService)
		userExamService service.UserExamService = service.NewUserExamService(userExamRepository)
		examService service.ExamService = service.NewExamService(examRepository, userExamService)

		// Controller
		userController controller.UserController = controller.NewUserController(userService)
		examController controller.ExamController = controller.NewExamController(examService)
		userExamController controller.UserExamController = controller.NewUserExamController(userExamService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// routes
	routes.UserRoutes(server, userController, jwtService)
	routes.ExamRoutes(server, examController, jwtService)
	routes.UserExamRoutes(server, userExamController, jwtService)

	server.Static("/assets", "./assets")
	port := os.Getenv("GOLANG_PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "localhost:" + port
	} else {
		serve = ":" + port
	}


	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
