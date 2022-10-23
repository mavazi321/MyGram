package controllers

import (
	"finalproject/database"
	"finalproject/helpers"
	"finalproject/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db := database.GetDB()

	err := db.Create(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"age": user.Age,
		"email":    user.Email,
		"id":       user.ID,
		"username": user.Username})
}

func UpdateUser(c *gin.Context) {
	var user models.User

	//get session user data
	userData := c.MustGet("userData").(jwt.MapClaims)

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	db := database.GetDB()
	errFoundEmail := db.Where("email=?", user.Email).Take(&models.User{}).Error

	if errFoundEmail == nil && userData["email"] != user.Email {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "email is taken"})
		return
	}

	errFoundUsername := db.Where("username=?", user.Username).Take(&models.User{}).Error

	if errFoundUsername == nil && userData["username"] != user.Username {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "username is taken"})
		return
	}

	username := user.Username
	email := user.Email

	res := db.Where("email=?", userData["email"]).First(&user)
	err := res.Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user.Username = username
	user.Email = email
	err1 := res.Updates(user).Error

	if err1 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"age":        user.Age,
		"updated_at": user.UpdatedAt})

}

func UserLogin(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	password := user.Password
	db := database.GetDB()
	err := db.Where("email = ?", user.Email).Take(&user).Error

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unatuhorized",
			"message": "invalid email/password",
		})
		return
	}

	comparePass := helpers.CheckPasswordHash(user.Password, password)

	if !comparePass {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "unatuhorized",
			"message": "invalid email/password",
		})
		return
	}

	token := helpers.GenerateToken(user.ID, user.Username, user.Email)

	c.JSON(http.StatusOK, gin.H{"token": token})

}

func DeleteUser(c *gin.Context) {
	var user models.User

	userData := c.MustGet("userData").(jwt.MapClaims)

	db := database.GetDB()

	err := db.Where("email = ?", userData["email"]).Take(&user).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	errDelete := db.Delete(&user).Error

	if errDelete != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errDelete,
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Your account has been succesfully deleted"})

}
