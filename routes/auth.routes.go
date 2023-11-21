package routes

import (
	"github.com/RamboXD/DB-Bonus/controllers"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
    authRouter := rg.Group("/auth")

    authRouter.POST("/login", rc.authController.SignInUser)

    registerRouter := authRouter.Group("/register")
    registerRouter.POST("/caregiver", rc.authController.SignUpCaregiver) 
    registerRouter.POST("/member", rc.authController.SignUpMember) 
}

