package models

import (
	"github.com/google/uuid"
)


type User struct {
    UserID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
    Email               string    `gorm:"type:varchar(255);unique;not null"`
    GivenName           string    `gorm:"type:varchar(100)"`
    Surname             string    `gorm:"type:varchar(100)"`
    City                string    `gorm:"type:varchar(100)"`
    PhoneNumber         string    `gorm:"type:varchar(15)"`
    ProfileDescription  string    `gorm:"type:text"`
    Password            string    `gorm:"type:varchar(255);not null"`
}