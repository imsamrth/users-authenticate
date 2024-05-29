package routes

import (
	"instix_auth/constants"
	constant "instix_auth/constants"

	"github.com/gin-gonic/gin"
)

func AuthlessRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Static("/my/profile-image", constant.ProfileImageDir)
	incomingRoutes.Static("/xchange/images", constant.ProductImageDir)
	incomingRoutes.Static("/assets", constant.AssetsDir)
	incomingRoutes.Static("/venue/images", constant.VenueImageDir)
	incomingRoutes.Static("body/logo", constants.BodyLogoDir)
}
