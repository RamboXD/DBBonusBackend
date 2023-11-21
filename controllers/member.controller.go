package controllers

import (
	"net/http"
	"time"

	"github.com/RamboXD/DB-Bonus/dto/request"
	"github.com/RamboXD/DB-Bonus/dto/response"
	"github.com/RamboXD/DB-Bonus/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type MemberController struct {
	DB *gorm.DB
}

func NewMemberController(DB *gorm.DB) MemberController {
	return MemberController{DB}
}

/*
Barlyk komekshilerdi tabu
=====================================================================================================================
*/

func (mc *MemberController) GetCaregivers(ctx *gin.Context) {
    var caregivers []models.Caregiver

    // Preload the User data for each Caregiver
    result := mc.DB.Preload("User").Find(&caregivers)
    if result.Error != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not retrieve caregivers"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success", "caregivers": caregivers})
}


func (mc *MemberController) CreateJob(ctx *gin.Context) {
    var jobReq request.JobCreationRequest

    if err := ctx.ShouldBindJSON(&jobReq); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        return
    }

    currentMember, exists := ctx.Get("currentMember")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Member not found"})
        return
    }
    member, ok := currentMember.(models.Member)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Member information is not valid"})
        return
    }

    // Create a new Job model from the request
    newJob := models.Job{
        MemberUserID:            member.MemberUserID,
        RequiredCaregivingType:  jobReq.RequiredCaregivingType,
        OtherRequirements:       jobReq.OtherRequirements,
        DatePosted:              time.Now(), 
    }

    // Save the new job in the database
    if err := mc.DB.Create(&newJob).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create job"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"status": "success", "job": newJob})
}

// Member sozdat etken jumystar
//=====================================================================================================================

func (mc *MemberController) MyJobs(ctx *gin.Context) {
    currentMember, exists := ctx.Get("currentMember")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Member not found"})
        return
    }
    member, ok := currentMember.(models.Member)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Member information is not valid"})
        return
    }

    var jobs []models.Job
    if err := mc.DB.Where("member_user_id = ?", member.MemberUserID).Find(&jobs).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve jobs"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success", "jobs": jobs})
}

// Jumysyna kyzykkan komekshilerdi karau
//=====================================================================================================================

func (mc *MemberController) JobApplicantsList(ctx *gin.Context) {
    jobIDStr := ctx.Param("id")
    jobID, err := uuid.Parse(jobIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid job ID"})
        return
    }

    currentMember, exists := ctx.Get("currentMember")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Member not found"})
        return
    }
    member, ok := currentMember.(models.Member)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Member information is not valid"})
        return
    }

    var job models.Job
    if err := mc.DB.Where("job_id = ? AND member_user_id = ?", jobID, member.MemberUserID).First(&job).Error; err != nil {
        ctx.JSON(http.StatusForbidden, gin.H{"status": "error", "message": "Access denied or job not found"})
        return
    }

    var applications []models.JobApplication
    if err := mc.DB.Where("job_id = ?", jobID).Find(&applications).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve job applications"})
        return
    }

	var caregivers []response.CaregiverDTO
    if err := mc.DB.Table("job_applications").
        Select("caregivers.*, users.email, users.given_name, users.surname, users.city, users.phone_number, users.profile_description").
        Joins("left join caregivers on caregivers.caregiver_user_id = job_applications.caregiver_user_id").
        Joins("left join users on users.user_id = caregivers.caregiver_user_id").
        Where("job_applications.job_id = ?", jobID).
        Scan(&caregivers).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve caregivers"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success", "caregivers": caregivers})
}


func (mc *MemberController) GetCaregiverDetails(ctx *gin.Context) {
    caregiverIDStr := ctx.Param("id")
    caregiverID, err := uuid.Parse(caregiverIDStr)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid caregiver ID"})
        return
    }

    var caregiverDetail response.CaregiverDTO
    if err := mc.DB.Table("caregivers").
        Select("caregivers.*, users.email, users.given_name, users.surname, users.city, users.phone_number, users.profile_description").
        Joins("left join users on users.user_id = caregivers.caregiver_user_id").
        Where("caregivers.caregiver_user_id = ?", caregiverID).
        Scan(&caregiverDetail).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Caregiver not found"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve caregiver details"})
        }
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success", "caregiver": caregiverDetail})
}

func (mc *MemberController) CreateAppointment(ctx *gin.Context) {
    var appointmentReq request.AppointmentCreationRequest

    if err := ctx.ShouldBindJSON(&appointmentReq); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        return
    }

    currentMember, exists := ctx.Get("currentMember")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Member not found"})
        return
    }
    member, ok := currentMember.(models.Member)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Member information is not valid"})
        return
    }

    // Check if the specified caregiver exists
    var caregiver models.Caregiver
    if err := mc.DB.First(&caregiver, "caregiver_user_id = ?", appointmentReq.CaregiverUserID).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Caregiver not found"})
        } else {
            ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to verify caregiver"})
        }
        return
    }

    newAppointment := models.Appointment{
        CaregiverUserID:   appointmentReq.CaregiverUserID,
        MemberUserID:      member.MemberUserID,
        AppointmentDate:   appointmentReq.AppointmentDate,
        AppointmentTime:   appointmentReq.AppointmentTime,
        WorkHours:         appointmentReq.WorkHours,
        Status:            "Pending", 
    }

    if err := mc.DB.Create(&newAppointment).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create appointment"})
        return
    }

    ctx.JSON(http.StatusCreated, gin.H{"status": "success", "appointment": newAppointment})
}

func (mc *MemberController) GetMemberAppointments(ctx *gin.Context) {
    currentMember, exists := ctx.Get("currentMember")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Member not found"})
        return
    }
    member, ok := currentMember.(models.Member)
    if !ok {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Member information is not valid"})
        return
    }

    var appointments []models.Appointment
    if err := mc.DB.Where("member_user_id = ?", member.MemberUserID).Preload("Caregiver.User").Find(&appointments).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to retrieve appointments"})
        return
    }

    ctx.JSON(http.StatusOK, gin.H{"status": "success", "appointments": appointments})
}