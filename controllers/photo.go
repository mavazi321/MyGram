package controllers

import (
	"finalproject/database"
	"finalproject/models"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreatePhoto(c *gin.Context) {
	var photo models.Photo
	var user models.User

	if err := c.ShouldBindJSON(&photo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//get user data
	userData := c.MustGet("userData").(jwt.MapClaims)

	db := database.GetDB()
	db.Where("email=?", userData["email"]).Take(&user)

	//set user data to photo
	photo.User = user
	photo.UserID = user.ID

	//create photo
	err := db.Create(&photo).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         photo.ID,
		"title":      photo.Title,
		"caption":    photo.Caption,
		"photo_url":  photo.PhotoUrl,
		"user_id":    photo.UserID,
		"created_at": photo.CreatedAt,
	})

}

func GetAllPhotos(c *gin.Context) {
	var photos []models.Photo

	//custome type to export
	type eksUser struct {
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	type eksPhoto struct {
		Id        uint      `json:"id"`
		Title     string    `json:"title"`
		Caption   string    `json:"caption"`
		PhotoUrl  string    `json:"photo_url"`
		UserId    uint      `json:"user_id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		User      eksUser   `json:"user"`
	}

	db := database.GetDB()

	err := db.Preload("User").Find(&photos).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var eksphotos []eksPhoto

	for _, val := range photos {
		//clone user
		var tuser eksUser
		tuser.Email = val.User.Email
		tuser.Username = val.User.Username

		var tphoto eksPhoto
		tphoto.Id = val.ID
		tphoto.Title = val.Title
		tphoto.Caption = val.Caption
		tphoto.UserId = val.UserID
		tphoto.CreatedAt = val.CreatedAt
		tphoto.UpdatedAt = val.UpdatedAt
		tphoto.User = tuser

		eksphotos = append(eksphotos, tphoto)

	}

	c.JSON(http.StatusOK, eksphotos)

	return

}

func UpdatePhotos(c *gin.Context) {
	var photoUpdated, photo models.Photo

	db := database.GetDB()
	err := db.Where("id=?", c.Param("photoId")).Take(&photoUpdated).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err1 := c.ShouldBindJSON(&photo); err1 != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}

	photoUpdated.Caption = photo.Caption
	photoUpdated.Title = photo.Title
	photoUpdated.PhotoUrl = photo.PhotoUrl

	errUpdate := db.Save(&photoUpdated).Error

	if errUpdate != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errUpdate.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": photoUpdated.ID,
		"title":      photoUpdated.Title,
		"caption":    photoUpdated.Caption,
		"photo_url":  photoUpdated.PhotoUrl,
		"user_id":    photoUpdated.UserID,
		"updated_at": photoUpdated.UpdatedAt})
}

func DeletePhotos(c *gin.Context) {
	var photo models.Photo

	db := database.GetDB()
	err := db.Where("id=?", c.Param("photoId")).Take(&photo).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "photos is not found"})
		return
	}

	errDelete := db.Delete(&photo).Error

	if errDelete != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errDelete.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "your photo has been succesfully deleted"})

}
