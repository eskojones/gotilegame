package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
)

func main() {
	g, err := makeGame("Game", 1024, 768, 1.0, 1024, "tileset1.png", 32)
	if err != nil {
		log.Fatal(err)
		return
	}
	// ebiten.SetFullscreen(true)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
