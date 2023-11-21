package routes

import (
	"github.com/RamboXD/DB-Bonus/controllers"
	"github.com/RamboXD/DB-Bonus/middleware"
	"github.com/gin-gonic/gin"
)

type MemberRouteController struct {
	memberController controllers.MemberController
}

func NewMemberRouteController(memberController controllers.MemberController) MemberRouteController {
	return MemberRouteController{memberController}
}

func (rc *MemberRouteController) MemberRoute(rg *gin.RouterGroup) {
    memberRouter := rg.Group("/member")

    memberRouter.GET("/getCaregivers", middleware.MemberCheck(), rc.memberController.GetCaregivers) 
    memberRouter.GET("/getAppointments", middleware.MemberCheck(), rc.memberController.GetMemberAppointments) 
    memberRouter.GET("/getCaregiver/:id", middleware.MemberCheck(), rc.memberController.GetCaregiverDetails) 
    memberRouter.POST("/createJob", middleware.MemberCheck(), rc.memberController.CreateJob) 
    memberRouter.POST("/createAppointment", middleware.MemberCheck(), rc.memberController.CreateAppointment) 
    memberRouter.GET("/jobs", middleware.MemberCheck(), rc.memberController.MyJobs) 
    memberRouter.GET("/job/applicants/:id", middleware.MemberCheck(), rc.memberController.JobApplicantsList) 
}

