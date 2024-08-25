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
	g.lastUpdate = time.Now().UnixMilli()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		os.Exit(0)
	}

	if g.player == nil {
		return nil
	}

	move := Point{X: g.player.position.X, Y: g.player.position.Y}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		move.Y -= g.player.moveSpeed * dt
		if move.Y < 0 {
			move.Y = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		move.Y += g.player.moveSpeed * dt
		if move.Y > float64(g.world.size) {
			move.Y = float64(g.world.size - 1)
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		move.X -= g.player.moveSpeed * dt
		if move.X < 0 {
			move.X = 0
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		move.X += g.player.moveSpeed * dt
		if move.X > float64(g.world.size) {
			move.X = float64(g.world.size - 1)
		}
	}

	g.SetEntityPosition(g.player, move.X, move.Y)

	return nil
}
