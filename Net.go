package main

import (
	"context"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func (net *NetConn) SendUpdate(g *Game) {
	if g.player == nil || time.Now().UnixMilli()-net.lastUpdate < 1000/PLAYER_UPDATE_PER_SECOND {
		return
	}
	net.lastUpdate = time.Now().UnixMilli()
	_, err := net.connection.Write([]byte(fmt.Sprintf("%s %.4f %.4f\n", CLIENT_FN_UPDATE, g.player.position.X, g.player.position.Y)))
	if err != nil {
		return
	}
}

func (net *NetConn) Update(g *Game) {
	for {
		if g.running == false {
			return
		}
		if net.connection == nil {
			// (re)connect
			ctx, cancel := context.WithTimeout(context.Background(), NET_TIMEOUT*time.Second)
			conn, err := net.dialer.DialContext(ctx, "tcp", g.net.server)
			if err != nil {
				log.Fatalf("Failed to connect: %v", err)
			}
			cancel()
			net.connection = conn
			time.Sleep(1 * time.Second)
			_, _ = conn.Write([]byte(fmt.Sprintf("%s %s %s\n", CLIENT_FN_CREATE, net.username, net.password)))
			_, _ = conn.Write([]byte(fmt.Sprintf("%s %s %s\n", CLIENT_FN_LOGIN, net.username, net.password)))
		}

		readBuf := make([]byte, NET_MSG_MAX_LEN)
		messageBuf := make([]byte, NET_MSG_MAX_LEN)
		var bytesReadCount int
		for {
			if g.running == false {
				return
			}
			_ = net.connection.SetReadDeadline(time.Now().Add(NET_READ_DEADLINE * time.Millisecond))
			count, err := net.connection.Read(readBuf)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) {
					break
				} else if !errors.Is(err, os.ErrDeadlineExceeded) {
					fmt.Printf("[read error: %s]\n", err)
					break
				}
			}

			// check if its time to send an update
			net.SendUpdate(g)

			if count == 0 {
				continue
			}

			messageBuf = fmt.Appendf(messageBuf[:bytesReadCount], "%s", readBuf[:count])
			bytesReadCount += count

			if bytesReadCount > NET_MSG_MAX_LEN {
				fmt.Printf("[%s sent an invalid message (too long)]\n", net.server)
				break
			}

			if strings.Contains(string(readBuf), "\n") {
				// split messageBuf by newline, and process each element as a separate message
				messages := strings.Split(string(messageBuf), "\n")
				for _, m := range messages {
					if len(m) > 0 {
						// fmt.Printf("server: %s\n", m)
						g.messages = append(g.messages, m)
					}
				}
				g.HandleMessages()
				clear(messageBuf)
				bytesReadCount = 0
			}
			clear(readBuf)
		}
	}
}

func (g *Game) HandleMessages() {
	for _, m := range g.messages {
		words := strings.Split(m, " ")
		if len(words) == 0 {
			continue
		}

		g.entityMutex.Lock()

		switch words[0] {
		case CLIENT_FN_UPDATE:
			// server is telling us about an entity's position
			id := words[1]
			xStr := words[2]
			yStr := words[3]
			x, _ := strconv.ParseFloat(xStr, 64)
			y, _ := strconv.ParseFloat(yStr, 64)
			// fmt.Printf("message: %s %s %.4f %.4f\n", CLIENT_FN_UPDATE, id, x, y)
			g.UpdateEntityPosition(id, x, y)
		case CLIENT_FN_QUERY:
			// server is responding to our query about an entity
			id := words[1]
			if g.world.entitiesFlat[id] != nil {
				spriteStr := words[2:]
				var frames []image.Point
				for _, xyStr := range spriteStr {
					if len(xyStr) == 0 {
						continue
					}
					parts := strings.Split(xyStr, ",")
					x, _ := strconv.ParseInt(parts[0], 10, 64)
					y, _ := strconv.ParseInt(parts[1], 10, 64)
					frames = append(frames, image.Point{X: int(x), Y: int(y)})
				}
				g.world.entitiesFlat[id].sprite = makeSprite(g.world.tileAtlas, frames, g.world.tileSize, 500)
			}
		}

		g.entityMutex.Unlock()

	}
	clear(g.messages)
}

func makeEntity(g *Game, id string) (*Entity, error) {
	if g.world.entitiesFlat[id] != nil {
		return nil, errors.New("entity already exists")
	}
	ent := new(Entity)
	ent.id = id
	ent.name = id
	ent.tileSize = g.world.tileSize
	ent.tileAtlas = g.world.tileAtlas
	ent.moveSpeed = 10.0
	g.world.entitiesFlat[id] = ent
	return ent, nil
}

func (g *Game) SetEntityPosition(e *Entity, x float64, y float64) bool {
	g.RemoveEntity(e)
	yMap := g.world.entities[int(y)]
	if yMap == nil {
		g.world.entities[int(y)] = make(map[int]map[string]*Entity)
		yMap = g.world.entities[int(y)]
	}
	xMap := yMap[int(x)]
	if xMap == nil {
		yMap[int(x)] = make(map[string]*Entity)
		xMap = yMap[int(x)]
	}
	if xMap[e.id] != nil {
		return false
	}
	xMap[e.id] = e
	e.position.X = x
	e.position.Y = y
	return true

}

func (g *Game) RemoveEntity(e *Entity) bool {
	yMap := g.world.entities[int(e.position.Y)]
	if yMap == nil {
		return false
	}
	xMap := yMap[int(e.position.X)]
	if xMap == nil {
		return false
	}
	if xMap[e.id] == nil {
		return false
	}
	delete(xMap, e.id)
	return true
}

func (g *Game) UpdateEntityPosition(id string, x float64, y float64) {
	if g.player != nil && g.player.id == id {
		// dont update ourselves from the network more than once
		return
	}
	entFlat := g.world.entitiesFlat[id]
	if entFlat == nil {
		// this is a new entity (to us)
		entNew, err := makeEntity(g, id)
		if err != nil {
			fmt.Println(err)
			return
		}
		entFlat = entNew
		if id == g.net.username {
			g.player = entFlat
		}

		// ask server about this entity (sprite, etc)
		g.SendEntityQuery(id)
	}
	g.SetEntityPosition(entFlat, x, y)
}

func (g *Game) GetEntitiesAt(x int, y int) []*Entity {
	yMap := g.world.entities[y]
	if yMap == nil {
		return nil
	}
	xMap := yMap[x]
	if xMap == nil {
		return nil
	}
	var result []*Entity
	for _, e := range xMap {
		result = append(result, e)
	}
	return result

}

func (g *Game) SendEntityQuery(id string) {
	if g.net.connection == nil || len(id) == 0 {
		return
	}
	_, err := g.net.connection.Write([]byte(fmt.Sprintf("%s %s\n", CLIENT_FN_QUERY, id)))
	if err != nil {
		return
	}
}
