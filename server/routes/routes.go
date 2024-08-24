package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/realjoni17/Hdocs/controllers"
)

func ConfigRoutes(router *gin.Engine) *gin.Engine {
	main := router.Group("v1")
	{
		user := main
		services := main
		cart := main
		payment := main
		{
			user.POST("/user", controllers.CreateUser)
			services.POST("/services", controllers.CreateService)
			services.GET("/services", controllers.GetService)
			cart.POST("/add/:service_id/:user_id", controllers.AddItemToCart)
			cart.GET("/cart", controllers.GetUserCart)
			payment.POST("/payment", controllers.Payment)
		}
	}
	return router
}
