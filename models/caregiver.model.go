package models

import (
	"github.com/google/uuid"
)


type Caregiver struct {
    CaregiverUserID   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Photo             []byte    `gorm:"type:bytea"`
    Gender            string    `gorm:"type:varchar(10)"`
    CaregivingType    string    `gorm:"type:varchar(100)"`
    HourlyRate        float64   `gorm:"type:numeric"`
    User              User      `gorm:"foreignKey:CaregiverUserID;references:UserID;constraint:OnDelete:CASCADE"`
}
