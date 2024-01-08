package main

import (
	"github.com/gin-gonic/gin"
	routes "github.com/mtoha/akhil/routes"
	"os"
)

func main() {
	port = os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	router := gin.new()
	router.User(gin.Logger())

	routes.AuthRoutes(router)
	routes.UserRoutes(router)

	router.Get("/API-1", func(c *gin.Context) {
		c.json(200, gin.H{"Success": "Access Granted for api-1"})
	})

	router.Get("/API-2", func(c *gin.Context) {
		c.json(200, gin.H{"Success": "Access Granted for api-2"})
	})

	router.Run(":" + port)
}
