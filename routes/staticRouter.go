package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func StaticRoutes(incomingRoutes *gin.Engine) {
	staticRoutes := incomingRoutes.Group("/static")
	{
		staticRoutes.PUT("/", controller.UploadSheet())
		staticRoutes.GET("/:sheet", controller.GetColumn())
	}
}
