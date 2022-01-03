package mmo

import (
	"gommo/engine/proceduralgeneration"
	"gommo/engine/tilemap"
	"math"
)

const (
	exponent       = 0.8
	islandExponent = 2.0
)

const (
	GrassTile tilemap.TileType = iota
	SandTile
	WaterTile
)

func CreateTilemap(seed int64, mapSize int, tileSize int) *tilemap.Tilemap {
	octaves := []proceduralgeneration.Octave{
		{Frequency: 0.02, Scale: 0.6},
		{Frequency: 0.05, Scale: 0.3},
		{Frequency: 0.1, Scale: 0.07},
		{Frequency: 0.2, Scale: 0.02},
		{Frequency: 0.4, Scale: 0.01},
	}
	terrain := proceduralgeneration.NewNoiseMap(seed, octaves, exponent)

	tiles := make([][]tilemap.Tile, mapSize)
	for x := range tiles {
		tiles[x] = make([]tilemap.Tile, mapSize)
		for y := range tiles[x] {
			height := terrain.Get(x, y)

			height = modifyHeightForIsland(mapSize, height, x, y)

			var tileType tilemap.TileType
			const waterLevel = 0.5
			const sandLevel = waterLevel + 0.1
			if height < waterLevel {
				tileType = WaterTile
			} else if height < sandLevel {
				tileType = SandTile
			} else {
				tileType = GrassTile
			}

			tiles[x][y] = tilemap.Tile{
				Type: tileType,
			}
		}
	}

	return tilemap.New(tiles, tileSize)
}

func modifyHeightForIsland(mapSize int, height float64, x int, y int) float64 {
	dx := float64(x)/float64(mapSize) - 0.5
	dy := float64(y)/float64(mapSize) - 0.5
	d := math.Sqrt(dx*dx+dy*dy) * 2
	d = math.Pow(d, islandExponent)
	height = (1 - d + height) / 2
	return height
}
