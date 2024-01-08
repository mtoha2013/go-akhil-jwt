package routes

import (
	"github.com/gin-gonic/gin"
	controller "github.com/mtoha/akhil/controllers"
	"github.com/mtoha/akhil/middleware"
)

func UserRoutes(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.POST("users/signup", controller.Signup())
	incomingRoutes.POST("users/login", controller.Login())
	incomingRoutes.GET("users", controller.GetUsers())
	incomingRoutes.GET("users/:user_id", controller.GetUsers())

}
