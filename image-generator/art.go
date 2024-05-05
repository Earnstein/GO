package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"sync"
	"time"

	"github.com/jdxyw/generativeart"
	"github.com/jdxyw/generativeart/common"
)

var (
	ColorSchema = []color.RGBA{
		{0xCF, 0x2B, 0x34, 0xFF},
		{0xF0, 0x8F, 0x46, 0xFF},
		{0xF0, 0xC1, 0x29, 0xFF},
		{0x19, 0x6E, 0x94, 0xFF},
		{0x35, 0x3A, 0x57, 0xFF},
	}
)

func DrawMany(drawings map[string]generativeart.Engine, wg *sync.WaitGroup) {
	wg.Add(len(drawings))
	for name := range drawings {
		go func(name string) {
			defer wg.Done()
			DrawOne(name)
		}(name)
	}
}

func DrawOne(name string) (string, error) {
	rand.New(rand.NewSource(time.Now().Unix()))
	engine, ok := DRAWINGS[name]
	if !ok {
		err := fmt.Errorf("engine %s not found", name)
		return "", err
	}
	
	c := generativeart.NewCanva(500, 500)
	c.SetBackground(color.RGBA{0xDF, 0xEB, 0xF5 , 0xFF})
	c.FillBackground()
	c.SetLineWidth(2.0)
	c.SetLineColor(common.Cheerful[rand.Intn(len(common.Cheerful))])
	c.SetColorSchema(common.Sleek)
	c.Draw(engine)
	filename := fmt.Sprintf("%s-art.png", name)
	err := c.ToPNG(filename)
	return filename, err
}
