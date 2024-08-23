package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"net"
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
	running       bool
	windowScale   float64
	worldSize     int
	tileSize      int
	tileset       *ebiten.Image
	locationTypes map[string]*LocationType
	world         map[int]map[int]*Location
	player        Player
	username      string
	password      string
	server        string
	dialer        net.Dialer
	connection    net.Conn
	readBuffer    [1024]byte
}

const CLIENT_FN_CREATE = "create" // command to create an account
const CLIENT_FN_LOGIN = "login"   // command to login to an account
const CLIENT_FN_LOGOUT = "logout" // command to logout from an account
const CLIENT_FN_UPDATE = "update" // command to update a player's position
