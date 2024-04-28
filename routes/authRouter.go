package routes

import (
	controller "instix_auth/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	incomingRoutes.GET("activate/:user_id", controller.ActivateGET())
	incomingRoutes.GET("activate/:user_id/:ver_hash", controller.EmailverGET())
}
