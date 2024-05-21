package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func VenueRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("venue/add", controller.CreateVenue())
}
