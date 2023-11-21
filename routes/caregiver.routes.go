package routes

import (
	"github.com/RamboXD/DB-Bonus/controllers"
	"github.com/RamboXD/DB-Bonus/middleware"
	"github.com/gin-gonic/gin"
)

type CaregiverRouteController struct {
	caregiverController controllers.CaregiverController
}

func NewCaregiverRouteController(caregiverController controllers.CaregiverController) CaregiverRouteController {
	return CaregiverRouteController{caregiverController}
}

func (rc *CaregiverRouteController) CaregiverRoute(rg *gin.RouterGroup) {
    caregiverRouter := rg.Group("/caregiver")

    caregiverRouter.GET("/getJobs", middleware.CaregiverCheck(), rc.caregiverController.GetJobs) 
    caregiverRouter.GET("/getAppointments", middleware.CaregiverCheck(), rc.caregiverController.GetAppointments) 
    caregiverRouter.POST("/applyJob/:id", middleware.CaregiverCheck(), rc.caregiverController.ApplyJob) 
    caregiverRouter.POST("/appointment/:id", middleware.CaregiverCheck(), rc.caregiverController.UpdateAppointmentStatus) 
}

