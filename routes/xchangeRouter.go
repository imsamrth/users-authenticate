package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func ItemRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("xchange/add", controller.CreateItem())
	incomingRoutes.GET("xchange/:product_id", controller.GetItem())
	incomingRoutes.GET("xchange", controller.GetItems())
	//sample query : "xchange?status=sent&date>2020-01-06T14:00:00.000Z&author.firstname=Jhon&skip=50&limit=100&sort=-date&fields=id,date"
	incomingRoutes.DELETE("xchange/:product_id", controller.DeleteItem())
	incomingRoutes.PATCH("xchange/:product_id", controller.UpdateItemInfo())
	incomingRoutes.PATCH("xchange/:product_id/images", controller.UpdateItemImages())
}
