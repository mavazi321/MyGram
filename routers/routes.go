package routers

import (
	"finalproject/controllers"
	"finalproject/middlewares"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	// Simple group: v1
	userRouter := router.Group("/users")
	{

		userRouter.POST("/register", controllers.CreateUser)

		userRouter.POST("/login", controllers.UserLogin)

		userRouter.PUT("/", middlewares.Authentication(), controllers.UpdateUser)

		userRouter.DELETE("/", middlewares.Authentication(), controllers.DeleteUser)
	}

	photosRouter := router.Group("/photos")
	{

		photosRouter.POST("/", middlewares.Authentication(), controllers.CreatePhoto)
		photosRouter.GET("/", middlewares.Authentication(), controllers.GetAllPhotos)
		photosRouter.PUT("/:photoId", middlewares.Authentication(), middlewares.PhotoAuthorization, controllers.UpdatePhotos)
		photosRouter.DELETE("/:photoId", middlewares.Authentication(), middlewares.PhotoAuthorization, controllers.DeletePhotos)

	}

	commentRouter := router.Group("/comments")
	{

		commentRouter.POST("/", middlewares.Authentication(), controllers.CreateComment)
		commentRouter.GET("/", middlewares.Authentication(), controllers.GetComments)
		commentRouter.PUT("/:commentId", middlewares.Authentication(), middlewares.CommentAuthorization, controllers.UpdateComment)
		commentRouter.DELETE("/:commentId", middlewares.Authentication(), middlewares.CommentAuthorization, controllers.DeleteComment)

	}

	socialMediaRouter := router.Group("/socialmedias")
	{

		socialMediaRouter.POST("/", middlewares.Authentication(), controllers.CreateSocialMedia)
		socialMediaRouter.GET("/", middlewares.Authentication(), controllers.GetSocialMedias)
		socialMediaRouter.PUT("/:socialMediaId", middlewares.Authentication(), middlewares.SocialMediaAuthorization, controllers.UpdateSocialMedia)
		socialMediaRouter.DELETE("/:socialMediaId", middlewares.Authentication(), middlewares.SocialMediaAuthorization, controllers.DeleteSocialMedia)

	}

	return router
}
