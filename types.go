package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"net"
	"sync"
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

// Entity : a dynamic entity within the game world
type Entity struct {
	id          string // hash?
	isPlayer    bool
	name        string
	position    Point   // position in the world
	netPosition Point   // position reported from server (lerp target)
	moveSpeed   float64 // speed that X/Y are allowed to change by per tick
	sprite      *Sprite
	tileSize    int
	tileAtlas   *ebiten.Image
}

// NetConn : all network-related vars
type NetConn struct {
	username   string     // username
	password   string     // password
	server     string     // host:port string
	dialer     net.Dialer // dialer responsible for establishing net.Conn
	connection net.Conn   // the tcp connection with server
	lastUpdate int64
}

type World struct {
	size          int                                // width / height of the game world (in tiles)
	tileSize      int                                // width / height of the tiles (in pixels)
	tileAtlas     *ebiten.Image                      // tile atlas to use for the location sprites
	sprites       map[string]*Sprite                 // sprites to use for drawing the locations
	locationTypes map[string]*LocationType           // archetypes defining locations
	locations     map[int]map[int]*Location          // the map itself
	entities      map[int]map[int]map[string]*Entity // game entities in the world
	entitiesFlat  map[string]*Entity                 // hash? -> Entity
}

type Game struct {
	running         bool     // is the game running?
	lastUpdate      int64    // unix timestamp in milliseconds of last update()
	windowWidth     int      //
	windowHeight    int      //
	windowScale     float64  // multiplier for the window resolution
	world           World    // all world-related vars
	player          *Entity  // the local player
	net             NetConn  // all network-related vars
	messages        []string //
	entityWaitGroup sync.WaitGroup
	entityMutex     sync.Mutex
}

const CLIENT_FN_CREATE = "create" // command to create an account
const CLIENT_FN_LOGIN = "login"   // command to login to an account
const CLIENT_FN_LOGOUT = "logout" // command to logout from an account
const CLIENT_FN_UPDATE = "update" // command to update a player's position
