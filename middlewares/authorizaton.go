package middlewares

import (
	"finalproject/database"
	"finalproject/models"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func SocialMediaAuthorization(c *gin.Context) {
	var socialMedia models.SocialMedia

	//get user data
	userData := c.MustGet("userData").(jwt.MapClaims)
	socialMediaParamId, _ := strconv.Atoi(c.Param("socialMediaId"))

	db := database.GetDB()
	errLoadingSocialMedia := db.Where("id=?", socialMediaParamId).Take(&socialMedia).Error

	if errLoadingSocialMedia != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if userData["id"] != float64(socialMedia.UserID) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}

func PhotoAuthorization(c *gin.Context) {
	var photo models.Photo

	photoId, _ := strconv.Atoi(c.Param("photoId"))

	db := database.GetDB()
	errLoading := db.Where("id=?", photoId).Take(&photo).Error

	if errLoading != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//get current user data
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := userData["id"]

	if float64(photo.UserID) != userId {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

}

func CommentAuthorization(c *gin.Context) {
	var comment models.Comment

	commentId, _ := strconv.Atoi(c.Param("commentId"))

	db := database.GetDB()
	errLoading := db.Where("id=?", commentId).Take(&comment).Error

	if errLoading != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	//get current user data
	userData := c.MustGet("userData").(jwt.MapClaims)
	userId := userData["id"]

	if float64(comment.UserID) != userId {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
