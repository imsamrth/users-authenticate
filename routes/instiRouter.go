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
		instiRoutes.POST("/council/:body_id", controller.CreateCouncil())
		instiRoutes.POST("/member/:body_id", controller.AddMember())
		instiRoutes.PUT("/member/:member_id", controller.UpdateMember())
		instiRoutes.DELETE("/member/:member_id", controller.DeleteMember())
		instiRoutes.Use(middleware.Authenticate())
		instiRoutes.GET("/", controller.GetBodies())
		instiRoutes.GET("/:body_id", controller.GetBody())
		instiRoutes.GET("/struct/:body_id", controller.GetCouncilStruct())
		instiRoutes.GET("/council/", controller.GetMembers())
		instiRoutes.GET("/member/:member_id", controller.GetMember())
		instiRoutes.PATCH("/approve/:body_id", controller.VerfiyBody())
	}
	eventRoutes := incomingRoutes.Group("/event")
	{
		eventRoutes.POST("/add", controller.CreateEvent())
		eventRoutes.PUT("/:event_id", controller.UpdateEvent())
		eventRoutes.PUT("/banner/:event_id", controller.PutEventBanner())
		eventRoutes.DELETE("/:event_id", controller.DeleteEvent())
		eventRoutes.Use(middleware.Authenticate())
		eventRoutes.GET("/", controller.GetEvents())
		eventRoutes.GET("/:event_id", controller.GetEvent())
		eventRoutes.GET("/banner/:event_id", controller.GetEventBanner())
	}
}
