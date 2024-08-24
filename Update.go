package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	"time"
)

func (g *Game) DeltaTime() float64 {
	return float64(time.Now().UnixMilli()-g.lastUpdate) * 0.001
}

func (g *Game) Update() error {
	dt := g.DeltaTime()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.position.Y -= g.player.moveSpeed * dt
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.position.Y += g.player.moveSpeed * dt
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.position.X -= g.player.moveSpeed * dt
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.position.X += g.player.moveSpeed * dt
	}

	g.lastUpdate = time.Now().UnixMilli()
	return nil
}
