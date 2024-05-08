package routes

import (
	controller "instix_auth/controllers"
	middleware "instix_auth/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	// TODO : uncomment the middleware Authenticate
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users/:user_id", controller.GetUser())
	incomingRoutes.GET("/users", controller.GetUsers())
}
