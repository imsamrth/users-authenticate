package routes

import (
	constant "instix_auth/constants"

	"github.com/gin-gonic/gin"
)

func AuthlessRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Static("/my/profile-image", constant.ProfileImageDir)
	incomingRoutes.Static("/assets", constant.AssetsDir)
}
