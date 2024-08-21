package main

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"math/rand/v2"
	"os"
)

func makeLocationType(g *Game, name string, tileXY []image.Point, isBlocking bool) *LocationType {
	locType := new(LocationType)
	locType.name = name
	locType.blocking = isBlocking
	tiles := make([]*ebiten.Image, len(tileXY))
	for i, p := range tileXY {
		tiles[i] = g.tileset.SubImage(image.Rect(p.X*g.tileSize, p.Y*g.tileSize, p.X*g.tileSize+g.tileSize, p.Y*g.tileSize+g.tileSize)).(*ebiten.Image)
	}
	locType.tiles = tiles
	g.locationTypes[name] = locType
	return locType
}

func makeWorld(g *Game) {
	// some test tiles
	makeLocationType(g, "grass", []image.Point{{8, 2}}, false)
	makeLocationType(g, "stone", []image.Point{{1, 2}}, false)
	makeLocationType(g, "stone2", []image.Point{{2, 2}}, false)
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
	g.player.x = 0 // rand.Float64() * float64(g.worldSize)
	g.player.y = 0 // rand.Float64() * float64(g.worldSize)
	g.player.moveSpeed = 0.0005
	g.player.tile = g.tileset.SubImage(image.Rect(1*g.tileSize, 0*g.tileSize, 1*g.tileSize+g.tileSize, 0*g.tileSize+g.tileSize)).(*ebiten.Image)
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
	return g, nil
}
