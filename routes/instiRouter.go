package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func InstiRoutes(incomingRoutes *gin.Engine) {
	instiRoutes := incomingRoutes.Group("/body")
	{
		instiRoutes.POST("/add", controller.CreateBody())
	}
}
