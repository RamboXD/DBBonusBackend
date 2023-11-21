package models

import (
	"time"

	"github.com/google/uuid"
)


type JobApplication struct {
    CaregiverUserID uuid.UUID `gorm:"type:uuid;primary_key"`
    JobID           uuid.UUID `gorm:"type:uuid;primary_key"`
    DateApplied     time.Time `gorm:"type:date"`
    Caregiver       Caregiver `gorm:"foreignKey:CaregiverUserID;references:CaregiverUserID;constraint:OnDelete:CASCADE"`
    Job             Job       `gorm:"foreignKey:JobID;references:JobID;constraint:OnDelete:CASCADE"`
}
