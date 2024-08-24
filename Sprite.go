package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
	"time"
)

func makeSprite(atlas *ebiten.Image, tiles []image.Point, size int, delay int64) *Sprite {
	sprite := new(Sprite)
	sprite.delay = delay
	sprite.lastFrame = time.Now().UnixMilli()
	sprite.frame = 0

	sprite.tiles = make([]*ebiten.Image, 0)
	for _, pos := range tiles {
		rect := image.Rect(pos.X*size, pos.Y*size, pos.X*size+size, pos.Y*size+size)
		tile := atlas.SubImage(rect).(*ebiten.Image)
		sprite.tiles = append(sprite.tiles, tile)
	}
	if len(sprite.tiles) > 1 {
		sprite.playing = true
	} else {
		sprite.playing = false
	}
	return sprite
}

func (sprite *Sprite) Draw(screen *ebiten.Image, options *ebiten.DrawImageOptions) {
	if sprite.playing {
		if time.Now().UnixMilli()-sprite.lastFrame >= sprite.delay {
			// adjust frame counter by how many delays have elapsed
			elapsed := time.Now().UnixMilli() - sprite.lastFrame
			for elapsed >= sprite.delay {
				sprite.frame++
				if sprite.frame >= len(sprite.tiles) {
					sprite.frame = 0
				}
				elapsed -= sprite.delay
			}
			// keep remainder of time elapsed
			sprite.lastFrame = time.Now().UnixMilli() - elapsed
		}
	}
	screen.DrawImage(sprite.tiles[sprite.frame], options)
}
