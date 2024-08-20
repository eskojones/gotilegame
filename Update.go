package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
)

func (g *Game) Update() error {
	dt := 1000.0 / ebiten.ActualFPS()
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.y -= g.player.moveSpeed * dt
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.y += g.player.moveSpeed * dt
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.x -= g.player.moveSpeed * dt
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.x += g.player.moveSpeed * dt
	}
	return nil
}
