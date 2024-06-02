package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(incomingRoutes *gin.Engine) {
	profileRoutes := incomingRoutes.Group("/profile")
	{
		profileRoutes.PUT("/my/avatar", controller.UploadAvatar())
		profileRoutes.PATCH("/my/profile/:isPrimary", controller.UpdateProfile())
		profileRoutes.PUT("my/username", controller.UpdateUsername())
		profileRoutes.GET("/", controller.GetProfiles())
	}
	postRoutes := incomingRoutes.Group("/post")
	{
		postRoutes.POST("/", controller.CreatePost())
	}
}
