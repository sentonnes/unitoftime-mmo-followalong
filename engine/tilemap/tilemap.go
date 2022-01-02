package tilemap

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type TileType uint8

type Tile struct {
	Type   TileType
	Sprite *pixel.Sprite
}

type Tilemap struct {
	TileSize int // In Pixels
	tiles    [][]Tile
	batch    *pixel.Batch
}

func New(tiles [][]Tile, batch *pixel.Batch, tileSize int) *Tilemap {
	return &Tilemap{tileSize, tiles, batch}
}

func (tilemap *Tilemap) Rebatch() {
	for x := range tilemap.tiles {
		for y := range tilemap.tiles[x] {
			tile := tilemap.tiles[x][y]
			position := pixel.V(float64(x*tilemap.TileSize), float64(y*tilemap.TileSize))

			matrix := pixel.IM.Moved(position)
			tile.Sprite.Draw(tilemap.batch, matrix)
		}
	}
}

func (tilemap *Tilemap) Draw(window *pixelgl.Window) {
	tilemap.batch.Draw(window)
}
