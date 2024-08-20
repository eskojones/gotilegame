package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

func drawWorld(screen *ebiten.Image, g *Game, drawPlayer bool) {
	bounds := screen.Bounds()
	tilesPerRow := math.Ceil(float64(bounds.Size().X) / float64(g.tileSize))
	tilesPerCol := math.Ceil(float64(bounds.Size().Y) / float64(g.tileSize))
	xStart := g.player.x - tilesPerRow/2
	yStart := g.player.y - tilesPerCol/2

	geo := ebiten.GeoM{}
	for ty := 0; ty < int(tilesPerCol); ty++ {
		worldY := int(yStart) + ty
		if worldY < 0 || worldY >= g.worldSize {
			continue
		}
		for tx := 0; tx < int(tilesPerRow); tx++ {
			worldX := int(xStart) + tx
			if worldX < 0 || worldX >= g.worldSize {
				continue
			}
			geo.Reset()
			geo.Translate(float64(tx*g.tileSize), float64(ty*g.tileSize))
			location := g.world[int(yStart)+ty][int(xStart)+tx]
			for _, tile := range location.locationType.tiles {
				screen.DrawImage(tile, &ebiten.DrawImageOptions{GeoM: geo})
			}
			if drawPlayer &&
				math.Floor(xStart+float64(tx)) == math.Floor(g.player.x) &&
				math.Floor(yStart+float64(ty)) == math.Floor(g.player.y) {
				screen.DrawImage(g.player.tile, &ebiten.DrawImageOptions{GeoM: geo})
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	drawWorld(screen, g, true)
	x, y := ebiten.CursorPosition()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("%.0f %.0f [%d %d] (%d,%d)", ebiten.ActualFPS(), ebiten.ActualTPS(), screen.Bounds().Size().X, screen.Bounds().Size().Y, x, y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return int(float64(outsideWidth) * g.windowScale), int(float64(outsideHeight) * g.windowScale)
}
