package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func NetUpdate(g *Game) {
	for {
		if g.running == false {
			return
		}
		if g.net.connection == nil {
			// (re)connect
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			conn, err := g.net.dialer.DialContext(ctx, "tcp", g.net.server)
			if err != nil {
				log.Fatalf("Failed to connect: %v", err)
			}
			cancel()
			g.net.connection = conn
			time.Sleep(1 * time.Second)
			_, _ = conn.Write([]byte(fmt.Sprintf("%s %s %s\n", CLIENT_FN_CREATE, g.net.username, g.net.password)))
			_, _ = conn.Write([]byte(fmt.Sprintf("%s %s %s\n", CLIENT_FN_LOGIN, g.net.username, g.net.password)))
		}

		readBuf := make([]byte, 1024)
		messageBuf := make([]byte, 1024)
		var bytesReadCount int
		for {
			if g.running == false {
				return
			}
			_ = g.net.connection.SetReadDeadline(time.Now().Add(50 * time.Millisecond))
			count, err := g.net.connection.Read(readBuf)
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
				fmt.Printf("[%s sent an invalid message (too long)]\n", g.net.server)
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
