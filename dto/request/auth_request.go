package request

import (
	"time"

	"github.com/RamboXD/DB-Bonus/models"
	"github.com/google/uuid"
)




type SignUpInput struct {
	Name            string `json:"name" binding:"required"`
	Email           string `json:"email" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	Photo           string `json:"photo" binding:"required"`
}

type SignInInput struct {
	Email    string `json:"email"  binding:"required"`
	Password string `json:"password"  binding:"required"`
}

type MemberSignUpInput struct {
    User   *models.User   `json:"user"`
    Member *models.Member `json:"member"`
    Address *models.Address `json:"address"`
}
type AppointmentCreationRequest struct {
    CaregiverUserID   uuid.UUID   `json:"caregiverUserId"`
    AppointmentDate   time.Time   `json:"appointmentDate"`
    AppointmentTime   time.Time   `json:"appointmentTime"`
    WorkHours         int         `json:"workHours"`
}

type CaregiverInfo struct {
    CaregiverUserID   uuid.UUID `json:"caregiverUserId"`
    Photo             []byte    `json:"photo"`
    Gender            string    `json:"gender"`
    CaregivingType    string    `json:"caregivingType"`
    HourlyRate        float64   `json:"hourlyRate"`
    Email             string    `json:"email"`
    GivenName         string    `json:"givenName"`
    Surname           string    `json:"surname"`
    City              string    `json:"city"`
    PhoneNumber       string    `json:"phoneNumber"`
    ProfileDescription string   `json:"profileDescription"`
    // Add other fields as necessary
}


type JobCreationRequest struct {
    RequiredCaregivingType  string    `json:"required_caregiving_type"`
    OtherRequirements       string    `json:"other_requirements"`
}
type CaregiverSignUpInput struct {
    User  *models.User  `json:"user"`
    Caregiver *models.Caregiver `json:"caregiver"`
}

type MaintenancePersonSignUpInput struct {
    User              *models.User             `json:"user"`
    // MaintenancePerson *models.MaintenancePerson `json:"maintenance_person"`
}

type FuelingPersonSignUpInput struct {
    User          *models.User          `json:"user"`
    // FuelingPerson *models.FuelingPerson `json:"fueling_person"`
}