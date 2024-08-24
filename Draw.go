package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

func drawWorld(screen *ebiten.Image, g *Game, drawPlayer bool) {
	bounds := screen.Bounds()
	tilesPerRow := float64(bounds.Size().X)/float64(g.world.tileSize) + 1
	tilesPerCol := float64(bounds.Size().Y)/float64(g.world.tileSize) + 1
	xStart := math.Floor(g.player.position.X - math.Floor(tilesPerRow/2))
	yStart := math.Floor(g.player.position.Y - math.Floor(tilesPerCol/2))

	options := &ebiten.DrawImageOptions{}
	// options.ColorScale.Scale(1, 1, 1, 1)
	for ty := 0; ty < int(tilesPerCol); ty++ {
		worldY := int(yStart) + ty
		if worldY < 0 || worldY >= g.world.size || g.world.locations[worldY] == nil {
			continue
		}
		for tx := 0; tx < int(tilesPerRow); tx++ {
			worldX := int(xStart) + tx
			if worldX < 0 || worldX >= g.world.size || g.world.locations[worldY][worldX] == nil {
				continue
			}
			options.GeoM.Reset()
			options.GeoM.Translate(float64(tx*g.world.tileSize), float64(ty*g.world.tileSize))
			location := g.world.locations[worldY][worldX]
			sprites := location.locationType.sprites
			for _, sprite := range sprites {
				sprite.Draw(screen, options)
			}
			if drawPlayer &&
				math.Floor(xStart+float64(tx)) == math.Floor(g.player.position.X) &&
				math.Floor(yStart+float64(ty)) == math.Floor(g.player.position.Y) {
				sprite := g.player.sprite
				// options.ColorScale.Scale(0.5, 0.5, 0.5, 1)
				sprite.Draw(screen, options)
				// options.ColorScale.Scale(1, 1, 1, 1)
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
