package middlewares

import (
	"errors"
	"finalproject/database"
	"finalproject/helpers"
	"finalproject/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		verifyToken, err := helpers.VerifyToken(c)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unatuhorized",
				"message": err.Error(),
			})
			return

		}

		db := database.GetDB()
		userData := verifyToken.(jwt.MapClaims)

		var user models.User
		err1 := db.Where("id = ?", userData["id"]).First(&user).Error
		if err1 != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "unatuhorized",
				"message": errors.New("sign in to proceed"),
			})
		}

		(verifyToken.(jwt.MapClaims))["email"] = user.Email
		(verifyToken.(jwt.MapClaims))["username"] = user.Username

		c.Set("userData", verifyToken)
		c.Next()
	}
}
