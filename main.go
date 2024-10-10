package main

import (
	"log"
	"os"

	"github.com/adieos/ilits-devsecops-deploy/cmd"
	"github.com/adieos/ilits-devsecops-deploy/config"
	"github.com/adieos/ilits-devsecops-deploy/controller"
	"github.com/adieos/ilits-devsecops-deploy/middleware"
	"github.com/adieos/ilits-devsecops-deploy/repository"
	"github.com/adieos/ilits-devsecops-deploy/routes"
	"github.com/adieos/ilits-devsecops-deploy/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := config.SetUpDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	if len(os.Args) > 1 {
		cmd.Commands(db)
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
	routes.User(server, userController, jwtService)

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong! website deployed successfully :)",
		})
	})

	server.Static("/assets", "./assets")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
