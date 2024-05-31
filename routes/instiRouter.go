package routes

import (
	controller "instix_auth/controllers"
	"instix_auth/middleware"

	"github.com/gin-gonic/gin"
)

func InstiRoutes(incomingRoutes *gin.Engine) {
	instiRoutes := incomingRoutes.Group("/body")
	{
		instiRoutes.POST("/add", controller.CreateBody())
		instiRoutes.PUT("/:body_id", controller.UpdateBody())
		instiRoutes.PUT("/struct/:body_id", controller.PutCouncilStruct())
		instiRoutes.GET("/struct/:body_id", controller.GetCouncilStruct())
		instiRoutes.POST("/council/:body_id", controller.CreateCouncil())
		instiRoutes.Use(middleware.Authenticate())
		instiRoutes.GET("/", controller.GetBodies())
		instiRoutes.GET("/:body_id", controller.GetBody())
		instiRoutes.GET("/council/", controller.GetMembers())
		instiRoutes.PATCH("/approve/:body_id", controller.VerfiyBody())
	}
}
