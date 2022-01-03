package main

//go:generate packer --input images --stats

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gommo/engine/asset"
	"gommo/engine/proceduralgeneration"
	"gommo/engine/render"
	"gommo/engine/tilemap"
	_ "image/png"
	"math"
	"os"
)

func main() {
	pixelgl.Run(runGame)
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	windowsMaxX      = 1024
	windowsMaxY      = 768
	windowsName      = "MMO"
	windowsVSync     = true
	windowsResizable = true
	purpleGemPng     = "purple.png"
	redGemPng        = "red.png"
	packedJson       = "packed.json"
	waterPng         = "water.png"
	sandPng          = "sand.png"
	grassPng         = "grass.png"
	tileSize         = 16
	mapSize          = 1000
	exponent         = 0.8
	islandExponent   = 2.0
)

var seed = int64(12345)

var load *asset.Load
var window *pixelgl.Window

func runGame() {
	setupGame()
	runGameLoop()
}

func runGameLoop() {
	spritesheet, err := load.Spritesheet(packedJson)
	check(err)
	tmap := createTileMap(spritesheet)
	spawnPoint := createSpawnPoint()
	people := createPeople(spritesheet, spawnPoint)
	camera, zoomSpeed := createCamera()
	gameLoop(camera, zoomSpeed, people, tmap)
}

func createSpawnPoint() pixel.Vec {
	return pixel.V(float64((tileSize*mapSize)/2), float64((tileSize*mapSize)/2))
}

func gameLoop(camera *render.Camera, zoomSpeed float64, people []Person, tmap *tilemap.Tilemap) {
	for !window.JustPressed(pixelgl.KeyEscape) {
		window.Clear(pixel.RGB(0, 0, 0))

		scroll := window.MouseScroll()
		if scroll.Y != 0 {
			camera.Zoom += zoomSpeed * scroll.Y
		}

		for i := range people {
			people[i].HandleInput(window)
		}

		camera.Position = people[0].Position
		camera.Update()

		window.SetMatrix(camera.Matrix())
		tmap.Draw(window)
		for i := range people {
			people[i].Draw(window)
		}
		window.SetMatrix(pixel.IM)

		window.Update()
	}
}

func createTileMap(spritesheet *asset.Spritesheet) *tilemap.Tilemap {
	octaves := []proceduralgeneration.Octave{
		{0.02, 0.6},
		{0.05, 0.3},
		{0.1, 0.07},
		{0.2, 0.02},
		{0.4, 0.01},
	}
	terrain := proceduralgeneration.NewNoiseMap(seed, octaves, exponent)

	tiles := make([][]tilemap.Tile, mapSize)
	for x := range tiles {
		tiles[x] = make([]tilemap.Tile, mapSize)
		for y := range tiles[x] {
			height := terrain.Get(x, y)

			height = modifyHeightForIsland(height, x, y)

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

			tiles[x][y] = GetTile(spritesheet, tileType)
		}
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, spritesheet.Picture())
	tmap := tilemap.New(tiles, batch, tileSize)
	tmap.Rebatch()
	return tmap
}

func modifyHeightForIsland(height float64, x int, y int) float64 {
	dx := float64(x)/float64(mapSize) - 0.5
	dy := float64(y)/float64(mapSize) - 0.5
	d := math.Sqrt(dx*dx+dy*dy) * 2
	d = math.Pow(d, islandExponent)
	height = (1 - d + height) / 2
	return height
}

func createCamera() (*render.Camera, float64) {
	camera := render.NewCamera(window, 0, 0)
	zoomSpeed := 0.1
	return camera, zoomSpeed
}

func createPeople(spritesheet *asset.Spritesheet, spawnPoint pixel.Vec) []Person {
	var purpleGemSprite, err = spritesheet.Get(purpleGemPng)
	check(err)
	redGemSprite, err := spritesheet.Get(redGemPng)
	check(err)

	people := make([]Person, 0)
	newPerson := NewPerson(purpleGemSprite, spawnPoint, ArrowKeybinds, 20)
	people = append(people, newPerson)
	newPerson = NewPerson(redGemSprite, spawnPoint, AWSDKeybinds, 4)
	people = append(people, newPerson)
	return people
}

func setupGame() {
	setupLoad()
	setupWindow()
}

func setupLoad() {
	load = asset.NewLoad(os.DirFS("./"))
}

func setupWindow() {
	cfg := getWindowsConfig()

	win, err := pixelgl.NewWindow(cfg)
	check(err)

	win.SetSmooth(false)
	window = win
}

func getWindowsConfig() pixelgl.WindowConfig {
	cfg := pixelgl.WindowConfig{
		Title:     windowsName,
		Bounds:    pixel.R(0, 0, windowsMaxX, windowsMaxY),
		VSync:     windowsVSync,
		Resizable: windowsResizable,
	}
	return cfg
}

const (
	GrassTile tilemap.TileType = iota
	SandTile
	WaterTile
)

func GetTile(spritesheet *asset.Spritesheet, tileType tilemap.TileType) tilemap.Tile {
	var spriteName string

	switch tileType {
	case GrassTile:
		spriteName = grassPng
	case SandTile:
		spriteName = sandPng
	case WaterTile:
		spriteName = waterPng
	default:
		panic("unknown TileType")
	}

	sprite, err := spritesheet.Get(spriteName)
	check(err)

	return tilemap.Tile{
		Type:   tileType,
		Sprite: sprite,
	}
}
