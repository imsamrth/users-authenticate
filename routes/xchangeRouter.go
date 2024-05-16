package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func ItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("xchange/add", controller.CreateItem())
	incomingRoutes.GET("xchange/:product_id", controller.GetItem())
	incomingRoutes.GET("xchange", controller.GetItems())
}