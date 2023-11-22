package controllers

import (
	"net/http"
	"time"

	"github.com/RamboXD/DB-Bonus/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type CaregiverController struct {
	DB *gorm.DB
}

func NewCaregiverController(DB *gorm.DB) CaregiverController {
	return CaregiverController{DB}
}

/*
Barlyk jumystardy tabu
=====================================================================================================================
*/

func (cc *CaregiverController) GetJobs(ctx *gin.Context) {
    var jobs []models.Job



    // Retrieve all jobs that the caregiver hasn't applied for, preloading Member and User data
    result := cc.DB.Preload("Member").Preload("Member.User").Find(&jobs)
    if result.Error != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve jobs"})
        return
    }

    // Return the list of jobs
    ctx.JSON(http.StatusOK, gin.H{"status": "success", "jobs": jobs})
}

func (cc *CaregiverController) GetCaregiverApplications(ctx *gin.Context) {
    currentCaregiver, exists := ctx.Get("currentCaregiver")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Caregiver not found"})
        return
    }
    caregiver, ok := currentCaregiver.(models.Caregiver)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Caregiver information is not valid"})
        return
    }

    var applications []models.JobApplication

    // Retrieve all applications made by the caregiver, preloading the Job, Member, and User details
    result := cc.DB.Preload("Job").Preload("Job.Member").Preload("Job.Member.User").Where("caregiver_user_id = ?", caregiver.CaregiverUserID).Find(&applications)
    if result.Error != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve applications"})
        return
    }

    // Create a slice to store the structured response
    var structuredApplications []map[string]interface{}

    for _, app := range applications {
        structuredApp := map[string]interface{}{
            "DateApplied": app.DateApplied,
            "JobDetails": map[string]interface{}{
                "JobID": app.Job.JobID,
                "MemberUserID": app.Job.MemberUserID,
                "RequiredCaregivingType": app.Job.RequiredCaregivingType,
                "OtherRequirements": app.Job.OtherRequirements,
                "DatePosted": app.Job.DatePosted,
                "MemberDetails": map[string]interface{}{
                    "MemberUserID": app.Job.Member.MemberUserID,
                    "HouseRules": app.Job.Member.HouseRules,
                    "UserDetails": map[string]interface{}{
                        "UserID": app.Job.Member.User.UserID,
                        "Email": app.Job.Member.User.Email,
                        "GivenName": app.Job.Member.User.GivenName,
                        "Surname": app.Job.Member.User.Surname,
                        "City": app.Job.Member.User.City,
                        "PhoneNumber": app.Job.Member.User.PhoneNumber,
                        "ProfileDescription": app.Job.Member.User.ProfileDescription,
                    },
                },
            },
        }
        structuredApplications = append(structuredApplications, structuredApp)
    }

    // Return the structured list of applications
    ctx.JSON(http.StatusOK, gin.H{"status": "success", "applications": structuredApplications})

}




func (cc *CaregiverController) ApplyJob(ctx *gin.Context) {
    jobIDStr := ctx.Param("id")
    jobID, err := uuid.Parse(jobIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid job ID"})
        return
    }

    currentCaregiver, exists := ctx.Get("currentCaregiver")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Caregiver not found"})
        return
    }
    caregiver, ok := currentCaregiver.(models.Caregiver)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Caregiver information is not valid"})
        return
    }

    jobApplication := models.JobApplication{
        CaregiverUserID: caregiver.CaregiverUserID,
        JobID:           jobID, 
        DateApplied:     time.Now(),
    }

    if err := cc.DB.Create(&jobApplication).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to apply for job"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": "Job application successful"})
}

func (cc *CaregiverController) GetAppointments(ctx *gin.Context) {
    currentCaregiver, exists := ctx.Get("currentCaregiver")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Caregiver not found"})
        return
    }
    caregiver, ok := currentCaregiver.(models.Caregiver)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Caregiver information is not valid"})
        return
    }

    var appointments []models.Appointment
    if err := cc.DB.Preload("Member").Preload("Member.User").Where("caregiver_user_id = ?", caregiver.CaregiverUserID).Find(&appointments).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve appointments"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success", "appointments": appointments})
}


func (cc *CaregiverController) UpdateAppointmentStatus(ctx *gin.Context) {
    appointmentIDStr := ctx.Param("id")
    appointmentID, err := uuid.Parse(appointmentIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid appointment ID"})
        return
    }

    var statusUpdate struct {
        Status string `json:"status"`
    }
    if err := ctx.ShouldBindJSON(&statusUpdate); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        return
    }

    currentCaregiver, exists := ctx.Get("currentCaregiver")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Caregiver not found"})
        return
    }
    caregiver, ok := currentCaregiver.(models.Caregiver)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Caregiver information is not valid"})
        return
    }

    var appointment models.Appointment
    if err := cc.DB.Where("appointment_id = ? AND caregiver_user_id = ?", appointmentID, caregiver.CaregiverUserID).First(&appointment).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Appointment not found or access denied"})
        return
    }

    appointment.Status = statusUpdate.Status
    if err := cc.DB.Save(&appointment).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to update appointment status"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Appointment status updated"})
}
