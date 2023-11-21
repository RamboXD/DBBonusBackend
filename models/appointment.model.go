package models

import (
	"time"

	"github.com/google/uuid"
)

type Appointment struct {
    AppointmentID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    CaregiverUserID uuid.UUID `gorm:"type:uuid"`
    MemberUserID    uuid.UUID `gorm:"type:uuid"`
    AppointmentDate time.Time `gorm:"type:date"`
    AppointmentTime time.Time `gorm:"type:time"`
    WorkHours       int       `gorm:"type:int"`
    Status          string    `gorm:"type:varchar(50)"`
    Caregiver       Caregiver `gorm:"foreignKey:CaregiverUserID;references:CaregiverUserID;constraint:OnDelete:CASCADE"`
    Member          Member    `gorm:"foreignKey:MemberUserID;references:MemberUserID;constraint:OnDelete:CASCADE"`
}
