package main

import (
	routes "instix_auth/routes"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access grated for api-1"})
	})

	routes.AuthlessRoutes(router)
	routes.AuthRoutes(router)
	routes.InstiRoutes(router)
	routes.UserRoutes(router)
	routes.ProfileRoutes(router)
	routes.ItemRoutes(router)
	routes.VenueRoutes(router)

	router.GET("/api-1", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access grated for api-1"})
	})

	router.GET("/api-2", func(c *gin.Context) {
		c.JSON(200, gin.H{"success": "Access grated for api-2"})
	})

	router.Run(":" + port)
}
