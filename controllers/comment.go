package controllers

import (
	"finalproject/database"
	"finalproject/models"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var comment models.Comment
	var user models.User
	var photo models.Photo

	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	db := database.GetDB()

	userData := c.MustGet("userData").(jwt.MapClaims)

	//set user
	db.Where("id=?", uint(userData["id"].(float64))).Take(&user)
	comment.UserID = user.ID
	comment.User = user
	//set photo
	db.Where("id=?", comment.PhotoID).Take(&photo)
	comment.Photo = photo

	fmt.Printf("tess %+v", comment)

	//insert to database
	err := db.Create(&comment).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": comment.ID,
		"message":    comment.Message,
		"photo_id":   comment.PhotoID,
		"user_id":    comment.UserID,
		"created_at": comment.CreatedAt})

}

func GetComments(c *gin.Context) {
	var comments []models.Comment

	//custome type to export
	type eksUser struct {
		Id       uint   `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
	}

	type eksPhoto struct {
		Id       uint   `json:"id"`
		Title    string `json:"title"`
		Caption  string `json:"caption"`
		PhotoUrl string `json:"photo_url"`
		UserId   uint   `json:"user_id"`
	}

	type eksComment struct {
		Id        uint      `json:"id"`
		Message   string    `json:"message"`
		PhotoId   uint      `json:"photo_id"`
		UserId    uint      `json:"user_id"`
		UpdatedAt time.Time `json:"updated_at"`
		CreatedAt time.Time `json:"created_at"`
		User      eksUser   `json:"user"`
		Photo     eksPhoto  `json:"photo"`
	}

	db := database.GetDB()
	err := db.Preload("Photo").Preload("User").Find(&comments).Error

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var ekscomments []eksComment

	for _, val := range comments {
		var tuser eksUser
		tuser.Id = val.User.ID
		tuser.Email = val.User.Email
		tuser.Username = val.User.Username

		var tphoto eksPhoto
		tphoto.Id = val.PhotoID
		tphoto.Title = val.Photo.Title
		tphoto.Caption = val.Photo.Caption
		tphoto.PhotoUrl = val.Photo.PhotoUrl
		tphoto.UserId = val.Photo.UserID

		var tcomment eksComment
		tcomment.Id = val.ID
		tcomment.Message = val.Message
		tcomment.PhotoId = val.PhotoID
		tcomment.UserId = val.UserID
		tcomment.UpdatedAt = val.UpdatedAt
		tcomment.CreatedAt = val.CreatedAt
		tcomment.User = tuser
		tcomment.Photo = tphoto

		ekscomments = append(ekscomments, tcomment)

	}

	c.JSON(http.StatusOK, ekscomments)
}

func UpdateComment(c *gin.Context) {
	var comment, updatedComment models.Comment

	db := database.GetDB()
	errCommentNotFound := db.Where("id=?", c.Param("commentId")).Take(&updatedComment).Error

	if errCommentNotFound != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	if errRequestContent := c.ShouldBindJSON(&comment); errRequestContent != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "request should in json format"})
		return
	}

	updatedComment.Message = comment.Message

	if errUpdateComment := db.Save(&updatedComment).Error; errUpdateComment != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": updatedComment.ID,
		"message":    updatedComment.Message,
		"updated_at": updatedComment.UpdatedAt})

}

func DeleteComment(c *gin.Context) {
	var comment models.Comment

	db := database.GetDB()
	errCommentNotFound := db.Where("id=?", c.Param("commentId")).Take(&comment).Error

	if errCommentNotFound != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "comment not found"})
		return
	}

	errDelete := db.Delete(&comment).Error

	if errDelete != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "your comment has been succesfully deleted"})
	return
}
