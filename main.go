package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 4 {
		log.Fatal("Usage: go run . <address>:<port> <username> <password>")
		return
	}
	g, err := makeGame("Game", 1024, 768, 1.0, 1024, "tileset1.png", 32)
	if err != nil {
		log.Fatal(err)
		return
	}

	g.net.server = os.Args[1]
	g.net.username = os.Args[2]
	g.net.password = os.Args[3]
	go NetUpdate(g)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
