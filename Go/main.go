package main

import (
	"errors"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Water Demo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		if errors.Is(err, terminated) {
			return
		}
		log.Fatal(err)
	}
	<-make(chan bool)
}

// https://www.fermyon.com/blog/optimizing-tinygo-wasm
// env GOOS=js GOARCH=wasm go build -o water.wasm
// env GOOS=js GOARCH=wasm tinygo build -o water.wasm
// wasm-opt -O water.wasm -o small_water.wasm