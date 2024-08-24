package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"net"
)

// LocationType : an archetype for a tile location
type LocationType struct {
	sprites  []*Sprite
	name     string // name to describe this location
	blocking bool   // does this location block entities from entering it?
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

type Sprite struct {
	tiles     []*ebiten.Image
	delay     int64
	lastFrame int64
	playing   bool
	frame     int
}

// Player : the local player
type Player struct {
	position  Point   // position in the world
	moveSpeed float64 // speed that X/Y are allowed to change by per tick
	sprite    *Sprite
	tileAtlas *ebiten.Image
}

// NetConn : all network-related vars
type NetConn struct {
	username   string     // username
	password   string     // password
	server     string     // host:port string
	dialer     net.Dialer // dialer responsible for establishing net.Conn
	connection net.Conn   // the tcp connection with server
}

type World struct {
	size          int                       // width / height of the game world (in tiles)
	tileSize      int                       // width / height of the tiles (in pixels)
	tileAtlas     *ebiten.Image             // tile atlas to use for the location sprites
	sprites       map[string]*Sprite        // sprites to use for drawing the locations
	locationTypes map[string]*LocationType  // archetypes defining locations
	locations     map[int]map[int]*Location // the map itself
}

type Game struct {
	running     bool    // is the game running?
	lastUpdate  int64   // unix timestamp in milliseconds of last update()
	windowScale float64 // multiplier for the window resolution
	world       World   // all world-related vars
	player      Player  // the local player
	net         NetConn // all network-related vars
}

const CLIENT_FN_CREATE = "create" // command to create an account
const CLIENT_FN_LOGIN = "login"   // command to login to an account
const CLIENT_FN_LOGOUT = "logout" // command to logout from an account
const CLIENT_FN_UPDATE = "update" // command to update a player's position
