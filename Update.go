package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

// NetHandler maintains a connection to the server
func NetHandler(g *Game) {
	for {
		if g.running == false {
			return
		}
		if g.connection == nil {
			// (re)connect
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			conn, err := g.dialer.DialContext(ctx, "tcp", g.server)
			if err != nil {
				log.Fatalf("Failed to connect: %v", err)
			}
			cancel()
			g.connection = conn
			time.Sleep(1 * time.Second)
			_, _ = conn.Write([]byte(fmt.Sprintf("create %s %s\n", g.username, g.password)))
			_, _ = conn.Write([]byte(fmt.Sprintf("login %s %s\n", g.username, g.password)))
		}

		readBuf := make([]byte, 1024)
		messageBuf := make([]byte, 1024)
		var bytesReadCount int
		for {
			if g.running == false {
				return
			}
			_ = g.connection.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			count, err := g.connection.Read(readBuf)
			if err != nil {
				if errors.Is(err, io.EOF) || errors.Is(err, os.ErrClosed) {
					break
				} else if !errors.Is(err, os.ErrDeadlineExceeded) {
					fmt.Printf("[read error: %s]\n", err)
					break
				}
			}

			if count == 0 {
				continue
			}

			messageBuf = fmt.Appendf(messageBuf[:bytesReadCount], "%s", readBuf[:count])
			bytesReadCount += count

			if bytesReadCount > 1024 {
				fmt.Printf("[%s sent an invalid message (too long)]\n", g.server)
				break
			}

			if strings.Contains(string(readBuf), "\n") {
				// split messageBuf by newline, and process each element as a separate message
				messages := strings.Split(string(messageBuf), "\n")
				for _, m := range messages {
					if len(m) > 0 {
						fmt.Printf("server: %s\n", m)
					}
				}
				clear(messageBuf)
				bytesReadCount = 0
			}
			clear(readBuf)
		}
	}
}

func (g *Game) Update() error {
	// dt := 1000.0 / ebiten.ActualFPS()
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		if g.connection != nil {
			_ = g.connection.Close()
		}
		os.Exit(0)
	} else if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.y -= (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.y += (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.x -= (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	} else if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.x += (g.player.moveSpeed / 1.0) * float64(g.worldSize)
	}

	return nil
}
