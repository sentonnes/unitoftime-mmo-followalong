package render

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gommo/engine/asset"
	"gommo/engine/tilemap"
)

type TilemapRender struct {
	spritesheet  *asset.Spritesheet
	batch        *pixel.Batch
	tileToSprite map[tilemap.TileType]*pixel.Sprite
}

func NewTilemapRender(spritesheet *asset.Spritesheet, tileToSprite map[tilemap.TileType]*pixel.Sprite) *TilemapRender {
	return &TilemapRender{
		spritesheet:  spritesheet,
		batch:        pixel.NewBatch(&pixel.TrianglesData{}, spritesheet.Picture()),
		tileToSprite: tileToSprite,
	}
}

func (tilemapRender TilemapRender) Clear() {
	tilemapRender.batch.Clear()
}

func (tilemapRender TilemapRender) Batch(tilemap *tilemap.Tilemap) {
	for x := 0; x < tilemap.Width(); x++ {
		for y := 0; y < tilemap.Height(); y++ {
			tile, ok := tilemap.Get(x, y)
			if !ok {
				continue
			}
			position := pixel.V(float64(x*tilemap.TileSize), float64(y*tilemap.TileSize))

			sprite, ok := tilemapRender.tileToSprite[tile.Type]
			if !ok {
				panic("unable to find TileType")
			}

			matrix := pixel.IM.Moved(position)
			sprite.Draw(tilemapRender.batch, matrix)
		}
	}
}

func (tilemapRender *TilemapRender) Draw(window *pixelgl.Window) {
	tilemapRender.batch.Draw(window)
}
