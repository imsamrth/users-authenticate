package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func ProfileRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.PUT("/my/profile-image", controller.UploadImagePOST())
	incomingRoutes.PATCH("/my/profile/:isPrimary", controller.UpdateProfile())
	incomingRoutes.PUT("/my/profile/username", controller.UpdateUsername())
	incomingRoutes.GET("/profiles", controller.GetProfiles())
}
