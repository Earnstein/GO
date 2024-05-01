package routes

import (
	"JobquestApi/controller"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	router.POST("/users/user", controller.HandleCreateUser())
	router.GET("/users/:user_id", controller.HandleGetUser())
	router.PUT("/users/:user_id", controller.HandleUpdateUser())
}