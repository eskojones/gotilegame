package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type LocationType struct {
	tiles    []*ebiten.Image
	name     string
	blocking bool
}

type Location struct {
	locationType *LocationType
	// ...
}

type Player struct {
	x         float64
	y         float64
	moveSpeed float64
	tile      *ebiten.Image
}

type Game struct {
	windowScale   float64
	worldSize     int
	tileSize      int
	tileset       *ebiten.Image
	locationTypes map[string]*LocationType
	world         map[int]map[int]*Location
	player        Player
}
