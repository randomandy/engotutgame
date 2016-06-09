package main

import (
	"engo.io/engo"
	"github.com/randomandy/ivo/scene"
)

func main() {
	opts := engo.RunOptions{
		Title:  "Ivo",
		Width:  500,
		Height: 500,
	}
	engo.Run(opts, &scene.MenuScene{})
}
