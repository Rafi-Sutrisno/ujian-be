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
		userRepository        repository.UserRepository        = repository.NewUserRepository(db)
		classRepository       repository.ClassRepository       = repository.NewClassRepository(db)
		userClassRepository   repository.UserClassRepository   = repository.NewUserClassRepository(db)
		examRepository        repository.ExamRepository        = repository.NewExamRepository(db)
		// userExamRepository    repository.UserExamRepository    = repository.NewUserExamRepository(db)
		examLangRepository    repository.ExamLangRepository    = repository.NewExamLangRepository(db)
		problemRepository     repository.ProblemRepository     = repository.NewProblemRepository(db)
		testCaseRepository    repository.TestCaseRepository    = repository.NewTestCaseRepository(db)
		submissionRepository  repository.SubmissionRepository  = repository.NewSubmissionRepository(db)
	
		// Service
		userService        service.UserService        = service.NewUserService(userRepository, jwtService)
		classService       service.ClassService       = service.NewClassService(classRepository)
		userClassService   service.UserClassService   = service.NewUserClassService(userClassRepository)
		// userExamService    service.UserExamService    = service.NewUserExamService(userExamRepository)
		examService        service.ExamService        = service.NewExamService(examRepository)
		examLangService    service.ExamLangService    = service.NewExamLangService(examLangRepository)
		problemService     service.ProblemService     = service.NewProblemService(problemRepository)
		testCaseService    service.TestCaseService    = service.NewTestCaseService(testCaseRepository)
		submissionService  service.SubmissionService  = service.NewSubmissionService(submissionRepository)
	
		// Controller
		userController        controller.UserController        = controller.NewUserController(userService)
		classController       controller.ClassController       = controller.NewClassController(classService)
		userClassController   controller.UserClassController   = controller.NewUserClassController(userClassService)
		examController        controller.ExamController        = controller.NewExamController(examService)
		// userExamController    controller.UserExamController    = controller.NewUserExamController(userExamService)
		examLangController    controller.ExamLangController    = controller.NewExamLangController(examLangService)
		problemController     controller.ProblemController     = controller.NewProblemController(problemService)
		testCaseController    controller.TestCaseController    = controller.NewTestCaseController(testCaseService)
		submissionController  controller.SubmissionController  = controller.NewSubmissionController(submissionService)
	)
	

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	// routes
	routes.UserRoutes(server, userController, jwtService)
	routes.ClassRoutes(server, classController, jwtService)
	routes.UserClassRoutes(server, userClassController, jwtService)
	routes.ExamRoutes(server, examController, jwtService)
	routes.ExamLangRoutes(server, examLangController, jwtService)
	routes.ProblemRoutes(server, problemController, jwtService)
	routes.TestCaseRoutes(server, testCaseController, jwtService)
	routes.SubmissionRoutes(server, submissionController, jwtService)

	// routes.UserExamRoutes(server, userExamController, jwtService)


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
