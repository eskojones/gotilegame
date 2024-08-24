package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand/v2"
	"os"
)

func makeLocationType(g *Game, name string, sprites []*Sprite, isBlocking bool) *LocationType {
	locType := new(LocationType)
	locType.name = name
	locType.blocking = isBlocking
	locType.sprites = sprites
	g.locationTypes[name] = locType
	return locType
}

func makeWorld(g *Game) {
	// some test tiles
	grass := makeSprite(g.tileset, []image.Point{{8, 2}, {6, 2}}, 32, 333)
	stone := makeSprite(g.tileset, []image.Point{{1, 2}, {2, 2}}, 32, 1000)
	stone2 := makeSprite(g.tileset, []image.Point{{3, 2}}, 32, 0)
	makeLocationType(g, "grass", []*Sprite{grass}, false)
	makeLocationType(g, "stone", []*Sprite{stone}, false)
	makeLocationType(g, "stone2", []*Sprite{stone2}, false)
	// test map fully populated
	types := []string{"grass", "stone", "stone2"}
	for y := 0; y < g.worldSize; y++ {
		if g.world[y] == nil {
			g.world[y] = make(map[int]*Location)
		}
		for x := 0; x < g.worldSize; x++ {
			g.world[y][x] = new(Location)
			g.world[y][x].locationType = g.locationTypes[types[rand.Int()%3]]
		}
	}
	// player config
	g.player.position.X = 0 // rand.Float64() * float64(g.worldSize)
	g.player.position.Y = 0 // rand.Float64() * float64(g.worldSize)
	g.player.moveSpeed = 0.0005
	g.player.sprite = makeSprite(g.tileset, []image.Point{{0, 0}, {0, 1}, {1, 1}, {1, 0}}, 32, 200)
}

func makeGame(windowTitle string, windowWidth int, windowHeight int, windowScale float64, worldSize int, tilesetFilename string, tileSize int) (*Game, error) {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetTPS(60)
	g := new(Game)
	g.windowScale = windowScale
	tilesetImgBytes, err := os.ReadFile(tilesetFilename)
	if err != nil {
		return nil, err
	}
	tilesetImage, _, err := image.Decode(bytes.NewReader(tilesetImgBytes))
	if err != nil {
		return nil, err
	}
	g.tileset = ebiten.NewImageFromImage(tilesetImage)
	g.locationTypes = make(map[string]*LocationType)
	g.world = make(map[int]map[int]*Location)
	g.worldSize = worldSize
	g.tileSize = tileSize
	makeWorld(g)
	g.running = true
	return g, nil
}
