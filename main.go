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

		// Service
		userService service.UserService = service.NewUserService(userRepository, jwtService)

		// Controller
		userController controller.UserController = controller.NewUserController(userService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// routes
	routes.Routes(server, userController, jwtService)

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
