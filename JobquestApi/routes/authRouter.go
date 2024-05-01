package routes

import (
	"JobquestApi/controller"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	router.POST("/users/signup", controller.Signup())
	router.POST("/users/login", controller.Login())
	router.GET("/users/logout", controller.Logout())
}