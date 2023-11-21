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

    // Retrieve all jobs from the database
    result := cc.DB.Find(&jobs)
    if result.Error != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve jobs"})
        return
    }

    // Return the list of jobs
    ctx.JSON(http.StatusOK, gin.H{"status": "success", "jobs": jobs})
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
    if err := cc.DB.Where("caregiver_user_id = ?", caregiver.CaregiverUserID).Find(&appointments).Error; err != nil {
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
