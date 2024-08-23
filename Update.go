package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
)

func (g *Game) Update() error {
	// dt := 1000.0 / ebiten.ActualFPS()
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.position.Y -= (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.position.Y += (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.position.X -= (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.position.X += (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	}

	return nil
}
