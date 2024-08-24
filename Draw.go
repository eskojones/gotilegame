package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

func drawWorld(screen *ebiten.Image, g *Game, pos Point) {
	options := &ebiten.DrawImageOptions{}
	w := float64(g.windowWidth)
	h := float64(g.windowHeight)
	// screen width in tiles
	viewWidth := math.Ceil(w / (float64(g.world.tileSize) * (1.0 / g.windowScale)))
	viewHeight := math.Ceil(h / (float64(g.world.tileSize) * (1.0 / g.windowScale)))
	// tile width in pixels
	tileWidth := math.Ceil(w / viewWidth)
	tileHeight := math.Ceil(h / viewHeight)
	// world coords of top left of screen
	viewX := pos.X - viewWidth*0.5
	viewY := pos.Y - viewHeight*0.5
	// screen width in tiles with the additional tile to avoid underdrawing bottom/right
	vW := viewWidth + 1
	vH := viewHeight + 1
	vX := math.Floor(viewX)
	vY := math.Floor(viewY)
	// fractional tile to offset all drawing for smooth scrolling
	fvX := math.Round((viewX - vX) * -tileWidth)
	fvY := math.Round((viewY - vY) * -tileHeight)

	var x float64 = 0
	var y float64 = 0
	// draw the world locations
	for y < vH {
		if g.world.locations[int(vY+y)] != nil && g.world.locations[int(vY+y)][int(vX+x)] != nil {
			loc := g.world.locations[int(vY+y)][int(vX+x)]
			scaleX, scaleY := g.Layout(int(fvX+(x*tileWidth)), int(fvY+(y*tileHeight)))
			options.GeoM.Reset()
			options.GeoM.Translate(float64(scaleX), float64(scaleY))
			for _, sprite := range loc.locationType.sprites {
				sprite.Draw(screen, options)
			}
		}
		x++
		if x == vW {
			x = 0.0
			y++
		}
	}

	// draw the entities above the ground (not yet implemented)
	// x = 0
	// y = 0
	// for y < vH {
	// 	if g.world.locations[int(vY+y)] != nil && g.world.locations[int(vY+y)][int(vX+x)] != nil {
	// 		loc := g.world.locations[int(vY+y)][int(vX+x)]
	// 		for _, entities := range loc.entities {
	// 			// draw all the entities
	// 		}
	// 	}
	// 	x++
	// 	if x == vW {
	// 		x = 0.0
	// 		y++
	// 	}
	// }

	// draw the players above all else (just draws local player for testing atm)
	x = 0
	y = 0
	for y < vH {
		if int(vX+x) == int(pos.X) && int(vY+y) == int(pos.Y) {
			fpX := (pos.X - math.Floor(pos.X)) * tileWidth
			fpY := (pos.Y - math.Floor(pos.Y)) * tileHeight
			scaleX, scaleY := g.Layout(int(fvX+fpX+(x*tileWidth)), int(fvY+fpY+(y*tileHeight)))
			options.GeoM.Reset()
			options.GeoM.Translate(float64(scaleX), float64(scaleY))
			g.player.sprite.Draw(screen, options)
		}
		x++
		if x == vW {
			x = 0.0
			y++
		}
	}

}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	x, y := ebiten.CursorPosition()
	drawWorld(screen, g, g.player.position)

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
