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
	incomingRoutes.GET("venue", controller.GetVenues())
	//sample query : "xchange?status=sent&date>2020-01-06T14:00:00.000Z&author.firstname=Jhon&skip=50&limit=100&sort=-date&fields=id,date"
}
