package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

func drawWorld(screen *ebiten.Image, g *Game, drawPlayer bool) {
	bounds := screen.Bounds()
	tilesPerRow := float64(bounds.Size().X)/float64(g.tileSize) + 1
	tilesPerCol := float64(bounds.Size().Y)/float64(g.tileSize) + 1
	xStart := math.Floor(g.player.position.X - math.Floor(tilesPerRow/2))
	yStart := math.Floor(g.player.position.Y - math.Floor(tilesPerCol/2))

	geo := ebiten.GeoM{}
	for ty := 0; ty < int(tilesPerCol); ty++ {
		worldY := int(yStart) + ty
		if worldY < 0 || worldY >= g.worldSize || g.world[worldY] == nil {
			continue
		}
		for tx := 0; tx < int(tilesPerRow); tx++ {
			worldX := int(xStart) + tx
			if worldX < 0 || worldX >= g.worldSize || g.world[worldY][worldX] == nil {
				continue
			}
			geo.Reset()
			geo.Translate(float64(tx*g.tileSize), float64(ty*g.tileSize))
			location := g.world[worldY][worldX]
			for _, tile := range location.locationType.tiles {
				screen.DrawImage(tile, &ebiten.DrawImageOptions{GeoM: geo})
			}
			if drawPlayer &&
				math.Floor(xStart+float64(tx)) == math.Floor(g.player.position.X) &&
				math.Floor(yStart+float64(ty)) == math.Floor(g.player.position.Y) {
				screen.DrawImage(g.player.tile, &ebiten.DrawImageOptions{GeoM: geo})
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	drawWorld(screen, g, true)
	x, y := ebiten.CursorPosition()

	ebitenutil.DebugPrint(screen, fmt.Sprintf(
		"%.0fFPS %.0fTPS\nRes: %dx%d\nCur: %d,%d\nPlayer: %.0f,%.0f\n",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
		screen.Bounds().Size().X, screen.Bounds().Size().Y,
		x, y,
		math.Floor(g.player.position.X), math.Floor(g.player.position.Y),
	))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float64(outsideWidth) * g.windowScale), int(float64(outsideHeight) * g.windowScale)
}
