package controllers

import (
	"finalproject/database"
	"finalproject/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateSocialMedia(c *gin.Context) {
	var (
		socialMedia models.SocialMedia
		user        models.User
	)

	if errBindJson := c.ShouldBindJSON(&socialMedia); errBindJson != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	//get current user session data
	userData := c.MustGet("userData").(jwt.MapClaims)

	db := database.GetDB()

	//get current user from databaes
	db.Where("id=?", userData["id"]).Take(&user)

	socialMedia.UserID = user.ID
	socialMedia.User = user

	errSave := db.Create(&socialMedia).Error

	if errSave != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": socialMedia.ID,
		"name":             socialMedia.Name,
		"social_media_url": socialMedia.SocialMediaUrl,
		"user_id":          socialMedia.UserID,
		"created_at":       socialMedia.CreatedAt})

}

func GetSocialMedias(c *gin.Context) {
	var socialMedias []models.SocialMedia

	//custome type to export
	type eksUser struct {
		Id       uint   `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	type eksSocialMedia struct {
		Id             uint      `json:"id"`
		Name           string    `json:"name"`
		SocialMediaUrl string    `json:"social_media_url"`
		UserId         uint      `json:"user_id"`
		CreatedAt      time.Time `json:"created_at"`
		UpdatedAt      time.Time `json:"updated_at"`
		User           eksUser   `json:"user"`
	}

	db := database.GetDB()
	errLoading := db.Preload("User").Find(&socialMedias).Error

	if errLoading != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var ekssocialmedias []eksSocialMedia

	for _, val := range socialMedias {
		var tuser eksUser
		tuser.Id = val.UserID
		tuser.Username = val.User.Username
		tuser.Email = val.User.Email

		var tsocial eksSocialMedia
		tsocial.Id = val.ID
		tsocial.Name = val.Name
		tsocial.SocialMediaUrl = val.SocialMediaUrl
		tsocial.UserId = val.UserID
		tsocial.CreatedAt = val.CreatedAt
		tsocial.UpdatedAt = val.UpdatedAt
		tsocial.User = tuser

		ekssocialmedias = append(ekssocialmedias, tsocial)

	}

	c.JSON(http.StatusOK, ekssocialmedias)

}

func UpdateSocialMedia(c *gin.Context) {
	var socialMedia, updatedSocialMedia models.SocialMedia

	db := database.GetDB()
	errSocialMediaNotfound := db.Where("id=?", c.Param("socialMediaId")).Take(&updatedSocialMedia).Error

	if errSocialMediaNotfound != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "social media not found"})
		return
	}

	if errBindJson := c.ShouldBindJSON(&socialMedia); errBindJson != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	updatedSocialMedia.Name = socialMedia.Name
	updatedSocialMedia.SocialMediaUrl = socialMedia.SocialMediaUrl

	errSaving := db.Save(&updatedSocialMedia).Error

	if errSaving != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": updatedSocialMedia.ID,
		"name":             updatedSocialMedia.Name,
		"social_media_url": updatedSocialMedia.SocialMediaUrl,
		"user_id":          updatedSocialMedia.UserID,
		"updated_at":       updatedSocialMedia.UpdatedAt})

}

func DeleteSocialMedia(c *gin.Context) {
	var socialMedia models.SocialMedia

	db := database.GetDB()
	errLoading := db.Where("id=?", c.Param("socialMediaId")).Take(&socialMedia).Error

	if errLoading != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "social media is not found"})
		return
	}

	errDeleting := db.Delete(&socialMedia).Error

	if errDeleting != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "social media has been succesfully deleted"})
}
