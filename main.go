package main

import (
	"engo.io/engo"
	"github.com/randomandy/ivo/scene"
)

func main() {
	opts := engo.RunOptions{
		Title:  "Ivo",
		Width:  800,
		Height: 800,
	}
	engo.Run(opts, &scene.DefaultScene{})
}
