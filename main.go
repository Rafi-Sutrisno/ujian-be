package main

import (
	"context"
	"log"
	"mods/application/service"
	"mods/command"
	"mods/config"
	domain "mods/domain/repository"
	"mods/infrastructure/auth"
	"mods/infrastructure/repository"
	"mods/interface/controller"
	"mods/interface/middleware"
	"mods/interface/routes"
	"os"
	"time"

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

func init() {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		panic(err)
	}
	time.Local = loc
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
		authRepository		  domain.AuthRepo			   = auth.NewAuthRepository(db)
		userRepository        domain.UserRepository        = repository.NewUserRepository(db)
		classRepository       domain.ClassRepository       = repository.NewClassRepository(db)
		userClassRepository   domain.UserClassRepository   = repository.NewUserClassRepository(db)
		examRepository        domain.ExamRepository        = repository.NewExamRepository(db)
		examLangRepository    domain.ExamLangRepository    = repository.NewExamLangRepository(db)
		problemRepository     domain.ProblemRepository     = repository.NewProblemRepository(db)
		testCaseRepository    domain.TestCaseRepository    = repository.NewTestCaseRepository(db)
		submissionRepository  domain.SubmissionRepository  = repository.NewSubmissionRepository(db)
		examSessionRepository domain.ExamSessionRepository = repository.NewExamSessionRepository(db)
		examProblemRepository domain.ExamProblemRepository = repository.NewExamProblemRepository(db)
		languageRepository    domain.LanguageRepository    = repository.NewLanguageRepository(db)
		userDraftRepository    domain.UserDraftRepository    = repository.NewUserDraftRepository(db)
	
		// Service
		userService        service.UserService        = service.NewUserService(userRepository, jwtService)
		classService       service.ClassService       = service.NewClassService(classRepository, userClassRepository)
		userClassService   service.UserClassService   = service.NewUserClassService(userClassRepository, userRepository)
		examService        service.ExamService        = service.NewExamService(examRepository, problemRepository, examProblemRepository, examLangRepository, authRepository)
		examLangService    service.ExamLangService    = service.NewExamLangService(examLangRepository)
		problemService     service.ProblemService     = service.NewProblemService(problemRepository, authRepository)
		testCaseService    service.TestCaseService    = service.NewTestCaseService(testCaseRepository)
		submissionService  service.SubmissionService  = service.NewSubmissionService(submissionRepository, testCaseRepository, languageRepository, problemRepository, authRepository)
		examSessionService  service.ExamSessionService  = service.NewExamSessionService(examSessionRepository, authRepository)
		examProblemService  service.ExamProblemService  = service.NewExamProblemService(examProblemRepository)
		languageService     service.LanguageService     = service.NewLanguageService(languageRepository)
		userDraftService     service.UserDraftService     = service.NewUserDraftService(userDraftRepository)
	
		// Controller
		userController        controller.UserController        = controller.NewUserController(userService)
		classController       controller.ClassController       = controller.NewClassController(classService)
		userClassController   controller.UserClassController   = controller.NewUserClassController(userClassService)
		examController        controller.ExamController        = controller.NewExamController(examService)
		examLangController    controller.ExamLangController    = controller.NewExamLangController(examLangService)
		problemController     controller.ProblemController     = controller.NewProblemController(problemService)
		testCaseController    controller.TestCaseController    = controller.NewTestCaseController(testCaseService)
		submissionController  controller.SubmissionController  = controller.NewSubmissionController(submissionService)
		examSessionController  controller.ExamSessionController  = controller.NewExamSessionController(examSessionService)
		examProblemController  controller.ExamProblemController  = controller.NewExamProblemController(examProblemService)
		languageController     controller.LanguageController     = controller.NewLanguageController(languageService)
		userDraftController   controller.UserDraftController     = controller.NewUserDraftController(userDraftService)
	)
	
	go submissionService.StartSubmissionPolling(context.Background())
	
	server := gin.Default()
	// server.Use(routes.RequestLogger())
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
	routes.ExamSessionRoutes(server, examSessionController, jwtService)
	routes.ExamProblemRoutes(server, examProblemController, jwtService)
	routes.LanguageRoutes(server, languageController, jwtService)
	routes.UserDraftRoutes(server, userDraftController, jwtService)

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
