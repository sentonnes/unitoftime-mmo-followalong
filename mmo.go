package mmo

import (
	"gommo/engine/ecs"
	"gommo/engine/physics"
	"gommo/engine/proceduralgeneration"
	"gommo/engine/tilemap"
	"math"
	"time"
)

const (
	exponent       = 0.8
	islandExponent = 2.0
)

const (
	GrassTile tilemap.TileType = iota
	SandTile
	WaterTile
	tileSize = 16
	mapSize  = 1000
)

var seed = int64(12345)

func LoadGame(engine *ecs.Engine) (*tilemap.Tilemap, ecs.Id, ecs.Id) {
	tmap := CreateTilemap(seed, mapSize, tileSize)

	spawnPoint := createSpawnPoint()
	purpleGemId := engine.NewId()
	ecs.Write(engine, purpleGemId, spawnPoint)
	ecs.Write(engine, purpleGemId, physics.Input{})

	redGemId := engine.NewId()
	ecs.Write(engine, redGemId, spawnPoint)
	ecs.Write(engine, redGemId, physics.Input{})

	return tmap, purpleGemId, redGemId
}

func createSpawnPoint() physics.Transform {
	return physics.Transform{X: float64((tileSize * mapSize) / 2), Y: float64((tileSize * mapSize) / 2)}
}

func CreateTilemap(seed int64, mapSize int, tileSize int) *tilemap.Tilemap {
	octaves := loadOctaves()
	terrain := proceduralgeneration.NewNoiseMap(seed, octaves, exponent)

	tiles := make([][]tilemap.Tile, mapSize)
	for x := range tiles {
		tiles[x] = make([]tilemap.Tile, mapSize)
		for y := range tiles[x] {
			height := terrain.Get(x, y)
			height = modifyHeightForIsland(mapSize, height, x, y)
			tileType := findTerrainTileType(height)
			tiles[x][y] = tilemap.Tile{Type: tileType}
		}
	}

	return tilemap.New(tiles, tileSize)
}

func loadOctaves() []proceduralgeneration.Octave {
	octaves := []proceduralgeneration.Octave{
		{Frequency: 0.02, Scale: 0.6},
		{Frequency: 0.05, Scale: 0.3},
		{Frequency: 0.1, Scale: 0.07},
		{Frequency: 0.2, Scale: 0.02},
		{Frequency: 0.4, Scale: 0.01},
	}
	return octaves
}

func findTerrainTileType(height float64) tilemap.TileType {
	const waterLevel = 0.5
	const sandLevel = waterLevel + 0.1
	var tileType tilemap.TileType
	if height < waterLevel {
		tileType = WaterTile
	} else if height < sandLevel {
		tileType = SandTile
	} else {
		tileType = GrassTile
	}
	return tileType
}

func modifyHeightForIsland(mapSize int, height float64, x int, y int) float64 {
	dx := float64(x)/float64(mapSize) - 0.5
	dy := float64(y)/float64(mapSize) - 0.5
	d := math.Sqrt(dx*dx+dy*dy) * 2
	d = math.Pow(d, islandExponent)
	height = (1 - d + height) / 2
	return height
}

func CreatePhysicsSystems(engine *ecs.Engine) []ecs.System {
	physicsSystems := []ecs.System{
		{Name: "HandleInput", Func: handleInputFunc(engine)},
	}
	return physicsSystems
}

func handleInputFunc(engine *ecs.Engine) func(dt time.Duration) {
	return func(dt time.Duration) {
		physics.HandleInput(engine)
	}
}
