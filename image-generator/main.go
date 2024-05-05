package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jdxyw/generativeart"
	"github.com/jdxyw/generativeart/arts"
)

var (
	DRAWINGS = map[string]generativeart.Engine{
		"maze": arts.NewMaze(10),
		"julia": arts.NewJulia(func(z complex128) complex128 {
			return z*z + complex(-0.001, 0.651)
		}, 40, 1.5, 1.5),
		"randcircle": arts.NewRandCicle(20, 80, 0.2, 2, 10, 40, true),
		"blackhole":  arts.NewBlackHole(200, 400, 0.01),
		"janus":      arts.NewJanus(5, 10),
		"random":     arts.NewRandomShape(150),
		"silksky":    arts.NewSilkSky(15, 5),
		"circles":    arts.NewColorCircle2(30),
	}
)

func startServer() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("templates/*.tmpl")
	imageRoutes := r.Group("/image")
	{
		imageRoutes.GET("/:name", HandleDrawImage)
	}


	listRoutes := r.Group("/list") 
	{
		listRoutes.GET("/", HandleEngineList)
	}
	return r
}


func handleHome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"MESSAGE": "TADA"})
}
func main() {

	// wg := sync.WaitGroup{}
	// DrawMany(DRAWINGS, &wg)
	// wg.Wait()
	r := startServer()
	r.GET("/", handleHome)
	r.Run(":8000")
}
