package models

import (
	"github.com/google/uuid"
)


type Address struct {
    MemberUserID uuid.UUID `gorm:"type:uuid;primary_key"`
    HouseNumber  string    `gorm:"type:varchar(10);primary_key"`
    Street       string    `gorm:"type:varchar(255);primary_key"`
    Town         string    `gorm:"type:varchar(100);primary_key"`
    Member       Member    `gorm:"foreignKey:MemberUserID;references:MemberUserID;constraint:OnDelete:CASCADE"`
}