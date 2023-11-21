package models

import (
	"time"

	"github.com/google/uuid"
)


type Job struct {
    JobID                   uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    MemberUserID            uuid.UUID `gorm:"type:uuid"`
    RequiredCaregivingType  string    `gorm:"type:varchar(100)"`
    OtherRequirements       string    `gorm:"type:text"`
    DatePosted              time.Time `gorm:"type:date"`
    Member                  Member    `gorm:"foreignKey:MemberUserID;references:MemberUserID;constraint:OnDelete:CASCADE"`
}
