package models

import (
	"github.com/google/uuid"
)


type Member struct {
    MemberUserID    uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    HouseRules      string    `gorm:"type:text"`
    User            User      `gorm:"foreignKey:MemberUserID;references:UserID;constraint:OnDelete:CASCADE"`
}
