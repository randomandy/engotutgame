package main

import (
	"engo.io/engo"
	"github.com/randomandy/ivo/scenes"
)

func main() {
	opts := engo.RunOptions{
		Title:  "Ivo",
		Width:  500,
		Height: 500,
	}
	engo.Run(opts, &scenes.MenuScene{})
}
