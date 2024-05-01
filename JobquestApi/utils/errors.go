package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)


type ErrorResponse struct {
	Message , Name ,Status string
}


const  (
	FAILED = "Failed"
)
func HandleInternalServerError(c *gin.Context, message string) {
	log.Println(message)
	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
		Status: FAILED,
		Message: message,
		Name: "Internal server error",
	} )
}

func HandleBadRequest(c *gin.Context, message string){
	log.Println(message)
	c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
		Status: FAILED,
		Message: message,
		Name: "BadRequest",
	})
}

func HandleUnauthenticated(c *gin.Context, message string){
	log.Println(message)
	c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
		Status: FAILED,
		Message: message,
		Name: "Unauthenticated",
	})
}

func HandleUnathourised(c *gin.Context, message string){
	log.Println(message)
	c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorResponse{
		Status: FAILED,
		Message: message,
		Name: "Unauthorised",
	})
}


func HandleNoContent(c *gin.Context, message string){
	log.Println(message)
	c.AbortWithStatusJSON(http.StatusNoContent, ErrorResponse{
		Status: FAILED,
		Message: message,
		Name: "No Content",
	})
}