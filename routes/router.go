package routes

import (
	"final-project/controllers"
	"final-project/middlewares"

	"github.com/gin-gonic/gin"
)

func StartApp() *gin.Engine {
	r := gin.Default()

	userRouter := r.Group("/users")
	{
		userRouter.POST("/register", controllers.UserRegister)

		userRouter.POST("/login", controllers.UserLogin)

		userRouter.Use(middlewares.Authentication())
		userRouter.PUT("/:userId", middlewares.UserAuthorization(), controllers.UpdateUser)
		userRouter.DELETE("/:userId", middlewares.UserAuthorization(), controllers.DeleteUser)
	}

	photoRouter := r.Group("/photos")
	{
		photoRouter.Use(middlewares.Authentication())
		photoRouter.POST("/", middlewares.CreateGetPhotoAuthorization(), controllers.CreatePhoto)
		
		photoRouter.GET("/", middlewares.CreateGetPhotoAuthorization(), controllers.GetPhoto)

		photoRouter.PUT("/:photoId", middlewares.UpdateDeletePhotoAuthorization(), controllers.UpdatePhoto)

		photoRouter.DELETE("/:photoId", middlewares.UpdateDeletePhotoAuthorization(), controllers.DeletePhoto)
	}

	commentRouter := r.Group("/comments")
	{
		commentRouter.Use(middlewares.Authentication())
		commentRouter.POST("/", middlewares.CreateGetCommentAuthorization(), controllers.CreateComment)
		
		commentRouter.GET("/", middlewares.CreateGetCommentAuthorization(), controllers.GetComent)

		commentRouter.PUT("/:commentId", middlewares.UpdateDeleteCommentAuthorization(), controllers.UpdateComment)

		commentRouter.DELETE("/:commentId", middlewares.UpdateDeleteCommentAuthorization(), controllers.DeleteComment)
	}

	socialMediaRouter := r.Group("/socialmedias")
	{
		socialMediaRouter.Use(middlewares.Authentication())
		socialMediaRouter.POST("/", middlewares.CreateGetSocialMediaAuthorization(), controllers.CreateSocialMedia)
		
		socialMediaRouter.GET("/", middlewares.CreateGetSocialMediaAuthorization(), controllers.GetSocialMedia)

		socialMediaRouter.PUT("/:socialMediaId", middlewares.UpdateDeleteSocialMediaAuthorization(), controllers.UpdateSocialMedia)

		socialMediaRouter.DELETE("/:socialMediaId", middlewares.UpdateDeleteSocialMediaAuthorization(), controllers.DeleteSocialMedia)
	}


	return r
}