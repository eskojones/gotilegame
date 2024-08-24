package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"os"
	"time"
)

func (g *Game) DeltaTimeUpdate() float64 {
	return float64(time.Now().UnixMilli()-g.lastUpdate) * 0.001
}

func (g *Game) Update() error {
	dt := g.DeltaTimeUpdate()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.position.Y -= g.player.moveSpeed * dt
		if g.player.position.Y < 0 {
			g.player.position.Y = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.position.Y += g.player.moveSpeed * dt
		if g.player.position.Y > float64(g.world.size) {
			g.player.position.Y = float64(g.world.size - 1)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.position.X -= g.player.moveSpeed * dt
		if g.player.position.X < 0 {
			g.player.position.X = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.position.X += g.player.moveSpeed * dt
		if g.player.position.X > float64(g.world.size) {
			g.player.position.X = float64(g.world.size - 1)
		}
	}

	g.lastUpdate = time.Now().UnixMilli()
	return nil
}
