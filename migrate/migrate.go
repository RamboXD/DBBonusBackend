package main

import (
	"fmt"
	"log"

	"github.com/RamboXD/DB-Bonus/initializers"
	"github.com/RamboXD/DB-Bonus/models"
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
}

func main() {
	initializers.DB.AutoMigrate(&models.User{})
	initializers.DB.AutoMigrate(&models.Caregiver{})
	initializers.DB.AutoMigrate(&models.Member{})
	initializers.DB.AutoMigrate(&models.Address{})
	initializers.DB.AutoMigrate(&models.Job{})
	initializers.DB.AutoMigrate(&models.JobApplication{})
	initializers.DB.AutoMigrate(&models.Appointment{})
	fmt.Println("? Migration complete")
}

