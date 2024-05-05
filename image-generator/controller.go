package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)


func HandleEngineList(c *gin.Context){
	c.HTML(http.StatusOK, "simple.tmpl", gin.H{
		"keys": maps.Keys(DRAWINGS),
	})
}

func HandleDrawImage(c *gin.Context) {
	name := c.Param("name")
	c.Header("Content-Type", "image/png")
	file, err := DrawOne(name)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"available engines": maps.Keys(DRAWINGS),
			"error":             err.Error()})
	}
	c.File(file)
}