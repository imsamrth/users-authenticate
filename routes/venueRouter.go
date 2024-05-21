package routes

import (
	"instix_auth/controllers"
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func VenueRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("venue/add", controller.CreateVenue())
	incomingRoutes.PUT("venue/:venue_id", controller.UpdateVenue())
	incomingRoutes.GET("venue/:venue_id", controllers.GetVenue())
}
