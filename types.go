package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"net"
)

// LocationType : an archetype for a tile location
type LocationType struct {
	tiles    []*ebiten.Image // all tiles to be drawn at this location (from 0th upwards)
	name     string          // name to describe this location
	blocking bool            // does this location block entities from entering it?
}

// Location : an instance of a tile location
type Location struct {
	locationType *LocationType // a reference to the archetype
	// ...
}

type Point struct {
	X float64
	Y float64
}

// Player : the local player
type Player struct {
	position  Point         // position in the world
	moveSpeed float64       // speed that X/Y are allowed to change by per tick
	tile      *ebiten.Image // todo: implement sprite system
}

// NetConn : all network-related vars
type NetConn struct {
	username   string     // username
	password   string     // password
	server     string     // host:port string
	dialer     net.Dialer // dialer responsible for establishing net.Conn
	connection net.Conn   // the tcp connection with server
}

type Game struct {
	running       bool                      // is the game running?
	windowScale   float64                   // multiplier for the window resolution
	worldSize     int                       // width / height of the game world (in tiles)
	tileSize      int                       // width / height of the tiles (in pixels)
	tileset       *ebiten.Image             // image file to use as tile atlas
	locationTypes map[string]*LocationType  // name -> location types
	world         map[int]map[int]*Location // y,x -> locations
	player        Player                    // the local player
	net           NetConn                   // all network-related vars
}

const CLIENT_FN_CREATE = "create" // command to create an account
const CLIENT_FN_LOGIN = "login"   // command to login to an account
const CLIENT_FN_LOGOUT = "logout" // command to logout from an account
const CLIENT_FN_UPDATE = "update" // command to update a player's position
