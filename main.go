package main

import (
	"log"
	"net/http"
	"time"

	"github.com/RamboXD/DB-Bonus/controllers"
	"github.com/RamboXD/DB-Bonus/initializers"
	"github.com/RamboXD/DB-Bonus/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	MemberController      controllers.MemberController
	MemberRouteController routes.MemberRouteController

	CaregiverController      controllers.CaregiverController
	CaregiverRouteController routes.CaregiverRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	MemberController = controllers.NewMemberController(initializers.DB)
	MemberRouteController = routes.NewMemberRouteController(MemberController)
	
	CaregiverController = controllers.NewCaregiverController(initializers.DB)
	CaregiverRouteController = routes.NewCaregiverRouteController(CaregiverController)
	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"https://d-bonus-front.vercel.app", "http://localhost:5173"}
    corsConfig.AllowMethods = []string{"POST", "GET", "PUT", "OPTIONS"}
    corsConfig.AllowHeaders = []string{"X-Requested-With","Origin", "Content-Type", "Authorization", "Accept", "User-Agent", "Cache-Control", "Pragma"}
    corsConfig.ExposeHeaders = []string{"Content-Length"}
    corsConfig.AllowCredentials = true
    corsConfig.MaxAge = 12 * time.Hour
	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		message := "Welcome to Golang with Gorm and Postgres"
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
	})

	AuthRouteController.AuthRoute(router)
	MemberRouteController.MemberRoute(router)
	CaregiverRouteController.CaregiverRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}


