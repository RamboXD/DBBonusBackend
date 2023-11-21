package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/RamboXD/DB-Bonus/dto/request"
	"github.com/RamboXD/DB-Bonus/dto/response"
	"github.com/RamboXD/DB-Bonus/initializers"
	"github.com/RamboXD/DB-Bonus/models"
	"github.com/RamboXD/DB-Bonus/utils"
	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(DB *gorm.DB) AuthController {
	return AuthController{DB}
}

/*
Caregiver registration
=====================================================================================================================
*/

func (ac *AuthController) SignUpCaregiver(ctx *gin.Context) {
    var payload request.CaregiverSignUpInput

    if err := ctx.ShouldBindJSON(&payload); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        return
    }

    if payload.User == nil || payload.Caregiver == nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User and Driver information are required"})
        return
    }

    hashedPassword, err := utils.HashPassword(payload.User.Password)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
        return
    }

    newUser := payload.User
    newUser.Email = strings.ToLower(newUser.Email)
    newUser.Password = hashedPassword

    result := ac.DB.Create(&newUser)
    if result.Error != nil {
        handleUserCreationError(ctx, result.Error)
        return
    }

    newCaregiver := payload.Caregiver
    newCaregiver.CaregiverUserID = newUser.UserID 

    caregiverResult := ac.DB.Create(&newCaregiver)
    if caregiverResult.Error != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to create caregiver profile"})
        return
    }

    userResponse := response.NewUserResponse(*newUser, newCaregiver, nil)
    ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}

/*
Member registration
=====================================================================================================================
*/

func (ac *AuthController) SignUpMember(ctx *gin.Context) {
    var payload request.MemberSignUpInput

    if err := ctx.ShouldBindJSON(&payload); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        return
    }

    if payload.User == nil || payload.Member == nil || payload.Address == nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "User and Driver information are required"})
        return
    }

    hashedPassword, err := utils.HashPassword(payload.User.Password)
    if err != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
        return
    }

    newUser := payload.User
    newUser.Email = strings.ToLower(newUser.Email)
    newUser.Password = hashedPassword

    result := ac.DB.Create(&newUser)
    if result.Error != nil {
        handleUserCreationError(ctx, result.Error)
        return
    }

    newMember := payload.Member
    newMember.MemberUserID = newUser.UserID 

    memberResult := ac.DB.Create(&newMember)
    if memberResult.Error != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to create caregiver profile"})
        return
    }

    newAddress := payload.Address
    newAddress.MemberUserID = newMember.MemberUserID 

    addressResult := ac.DB.Create(&newAddress)
    if addressResult.Error != nil {
        ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Failed to create caregiver profile"})
        return
    }

    userResponse := response.NewUserResponse(*newUser, nil, newMember)
    ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": gin.H{"user": userResponse}})
}



/*
Defaul Login HUMANITY
=====================================================================================================================
*/

func (ac *AuthController) SignInUser(ctx *gin.Context) {
    var payload *request.SignInInput
    
    if err := ctx.ShouldBindJSON(&payload); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        return
    }
    log.Printf(payload.Email)
    log.Printf(payload.Password)
    var user models.User
    result := ac.DB.First(&user, "email = ?", strings.ToLower(payload.Email))
    if result.Error != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
        return
    }
    
    if err := utils.VerifyPassword(user.Password, payload.Password); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Invalid email or Password"})
        return
    }

    // Determine the user's role
    var role string
    var member models.Member
    var caregiver models.Caregiver

    if err := ac.DB.First(&member, "member_user_id = ?", user.UserID).Error; err == nil {
        role = "member"
    } else if err := ac.DB.First(&caregiver, "caregiver_user_id = ?", user.UserID).Error; err == nil {
        role = "caregiver"
    } else {
        role = "unknown"
    }

    config, _ := initializers.LoadConfig(".")
    
    access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.UserID, config.AccessTokenPrivateKey)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": access_token, "role": role})
}


func handleUserCreationError(ctx *gin.Context, err error) {
	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": "User with that email already exists"})
	} else {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Something bad happened"})
	}
}