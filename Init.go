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
	g.world.locationTypes[name] = locType
	return locType
}

func makeWorld(g *Game) {
	// some test tiles (single frame sprites)
	dirt := makeSprite(g.world.tileAtlas, []image.Point{{11, 2}}, 32, 0)
	stone := makeSprite(g.world.tileAtlas, []image.Point{{11, 4}}, 32, 0)
	broken_stone := makeSprite(g.world.tileAtlas, []image.Point{{11, 5}}, 32, 0)
	// test locations composed of layered sprites
	makeLocationType(g, "stone over dirt", []*Sprite{dirt, broken_stone}, false)
	makeLocationType(g, "stone", []*Sprite{stone}, false)
	makeLocationType(g, "dirt", []*Sprite{dirt}, false)
	// test map fully populated
	for y := 0; y < g.world.size; y++ {
		if g.world.locations[y] == nil {
			g.world.locations[y] = make(map[int]*Location)
		}
		for x := 0; x < g.world.size; x++ {
			r := rand.Int() % 100
			if r > 95 {
				g.world.locations[y][x] = new(Location)
				g.world.locations[y][x].locationType = g.world.locationTypes["stone over dirt"]
			} else if r > 20 {
				g.world.locations[y][x] = new(Location)
				g.world.locations[y][x].locationType = g.world.locationTypes["stone"]
			}
		}
	}
	// player config
	// g.player = new(Entity)
	// g.player.isPlayer = true
	// g.player.name = "me"
	// g.player.position.X = 0 // rand.Float64() * float64(g.worldSize)
	// g.player.position.Y = 0 // rand.Float64() * float64(g.worldSize)
	// g.player.moveSpeed = 10.0
	// g.player.tileAtlas = g.world.tileAtlas
	// g.player.tileSize = g.world.tileSize
	// g.player.sprite = makeSprite(g.player.tileAtlas, []image.Point{{2, 14}, {2, 13}, {3, 13}}, 32, 200)

	g.world.entities = make(map[int]map[int]map[string]*Entity)
	g.world.entitiesFlat = make(map[string]*Entity)
}

func makeGame(windowTitle string, windowWidth int, windowHeight int, windowScale float64, worldSize int, tilesetFilename string, tileSize int) (*Game, error) {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle(windowTitle)
	ebiten.SetTPS(60)
	g := new(Game)
	g.windowWidth = windowWidth
	g.windowHeight = windowHeight
	g.windowScale = windowScale
	tilesetImgBytes, err := os.ReadFile(tilesetFilename)
	if err != nil {
		return nil, err
	}
	tilesetImage, _, err := image.Decode(bytes.NewReader(tilesetImgBytes))
	if err != nil {
		return nil, err
	}
	g.world.tileAtlas = ebiten.NewImageFromImage(tilesetImage)
	g.world.tileSize = tileSize
	g.world.size = worldSize
	g.world.locationTypes = make(map[string]*LocationType)
	g.world.locations = make(map[int]map[int]*Location)
	g.world.sprites = make(map[string]*Sprite)
	makeWorld(g)
	g.running = true
	return g, nil
}
