package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"math"
)

func (world *World) Draw(screen *ebiten.Image, g *Game, pos Point) {
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
	fvX := math.Floor((viewX - vX) * -tileWidth)
	fvY := math.Floor((viewY - vY) * -tileHeight)

	// draw the world locations
	for y := 0.0; y < vH; y++ {
		yMap := world.locations[int(vY+y)]
		if yMap == nil {
			continue
		}
		for x := 0.0; x < vW; x++ {
			loc := yMap[int(vX+x)]
			if loc == nil {
				continue
			}
			scaleX, scaleY := g.Layout(int(fvX+(x*tileWidth)), int(fvY+(y*tileHeight)))
			options.GeoM.Reset()
			options.GeoM.Translate(float64(scaleX), float64(scaleY))
			sprites := loc.locationType.sprites
			for _, sprite := range sprites {
				sprite.Draw(screen, options)
			}
		}
	}

	// draw the players above all else
	for y := 0.0; y < vH; y++ {
		yInt := int(vY + y)

		for x := 0.0; x < vW; x++ {
			xInt := int(vX + x)

			g.entityMutex.Lock()
			ents := g.GetEntitiesAt(xInt, yInt)
			g.entityMutex.Unlock()

			for _, ent := range ents {
				if ent.sprite == nil {
					continue
				}
				fpX := (ent.position.X - math.Floor(ent.position.X)) * tileWidth
				fpY := (ent.position.Y - math.Floor(ent.position.Y)) * tileHeight
				scaleX, scaleY := g.Layout(int(fvX+fpX+(x*tileWidth)), int(fvY+fpY+(y*tileHeight)))
				options.GeoM.Reset()
				options.GeoM.Translate(float64(scaleX), float64(scaleY))
				ent.sprite.Draw(screen, options)
			}
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Clear()
	if g.player == nil {
		return
	}

	x, y := ebiten.CursorPosition()
	g.world.Draw(screen, g, g.player.position)

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
