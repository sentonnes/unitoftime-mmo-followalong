package main

//go:generate packer --input images --stats

import (
	"context"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	mmo "gommo"
	"gommo/engine/asset"
	"gommo/engine/ecs"
	"gommo/engine/render"
	"gommo/engine/tilemap"
	_ "image/png"
	"log"
	"nhooyr.io/websocket"
	"os"
	"time"
)

func main() {
	url := "ws://localhost:8000"
	ctx := context.Background()
	c, resp, err := websocket.Dial(ctx, url, nil)
	check(err)

	log.Println("Connection Response:", resp)

	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)

	go func() {
		counter := byte(0)
		for {
			time.Sleep(1 * time.Second)
			n, err := conn.Write([]byte{counter})
			if err != nil {
				log.Println("error sending:", err)
				return
			}

			log.Println("sent n bytes: ", n)
			counter++
		}
	}()

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
)

var seed = int64(12345)

var load *asset.Load
var window *pixelgl.Window
var engine *ecs.Engine
var playerId ecs.Id

func runGame() {
	setupGame()
	runGameLoop()
}

func runGameLoop() {
	spritesheet, err := load.Spritesheet(packedJson)
	check(err)
	tmap := createTileMap(spritesheet)
	spawnPoint := createSpawnPoint()
	createPeople(spritesheet, spawnPoint)
	camera, zoomSpeed := createCamera()
	gameLoop(camera, zoomSpeed, tmap)
}

func createSpawnPoint() Transform {
	return Transform{float64((tileSize * mapSize) / 2), float64((tileSize * mapSize) / 2)}
}

func gameLoop(camera *render.Camera, zoomSpeed float64, tmapRender *render.TilemapRender) {
	for !window.JustPressed(pixelgl.KeyEscape) {
		window.Clear(pixel.RGB(0, 0, 0))

		scroll := window.MouseScroll()
		if scroll.Y != 0 {
			camera.Zoom += zoomSpeed * scroll.Y
		}

		HandleInput(window, engine)

		transform := Transform{}
		ok := ecs.Read(engine, playerId, &transform)
		if ok {
			camera.Position = pixel.V(transform.X, transform.Y)
		}
		camera.Update()

		window.SetMatrix(camera.Matrix())
		tmapRender.Draw(window)

		DrawSprite(window, engine)

		window.SetMatrix(pixel.IM)

		window.Update()
	}
}

func createTileMap(spritesheet *asset.Spritesheet) *render.TilemapRender {
	grassTile, err := spritesheet.Get(grassPng)
	check(err)
	sandTile, err := spritesheet.Get(sandPng)
	check(err)
	waterTile, err := spritesheet.Get(waterPng)
	check(err)

	tmap := mmo.CreateTilemap(seed, mapSize, tileSize)
	tmapRender := render.NewTilemapRender(spritesheet, map[tilemap.TileType]*pixel.Sprite{
		mmo.GrassTile: grassTile,
		mmo.SandTile:  sandTile,
		mmo.WaterTile: waterTile,
	})
	tmapRender.Batch(tmap)
	return tmapRender
}

func createCamera() (*render.Camera, float64) {
	camera := render.NewCamera(window, 0, 0)
	zoomSpeed := 0.1
	return camera, zoomSpeed
}

func createPeople(spritesheet *asset.Spritesheet, spawnPoint Transform) {
	var purpleGemSprite, err = spritesheet.Get(purpleGemPng)
	check(err)
	redGemSprite, err := spritesheet.Get(redGemPng)
	check(err)

	purpleGemId := engine.NewId()
	ecs.Write(engine, purpleGemId, Sprite{purpleGemSprite})
	ecs.Write(engine, purpleGemId, spawnPoint)
	ecs.Write(engine, purpleGemId, AWSDKeybinds)

	playerId = purpleGemId

	redGemId := engine.NewId()
	ecs.Write(engine, redGemId, Sprite{redGemSprite})
	ecs.Write(engine, redGemId, spawnPoint)
	ecs.Write(engine, redGemId, ArrowKeybinds)
}

func setupGame() {
	setupLoad()
	setupEngine()
	setupWindow()
}

func setupEngine() {
	engine = ecs.NewEngine()
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

func DrawSprite(window *pixelgl.Window, engine *ecs.Engine) {
	ecs.Each(engine, Sprite{}, func(id ecs.Id, a interface{}) {
		sprite := a.(Sprite)

		transform := Transform{}
		ok := ecs.Read(engine, id, &transform)
		if !ok {
			return
		}

		position := pixel.V(transform.X, transform.Y)
		sprite.Draw(window, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(position))
	})
}

func HandleInput(window *pixelgl.Window, engine *ecs.Engine) {
	ecs.Each(engine, Keybinds{}, func(id ecs.Id, a interface{}) {
		keybinds := a.(Keybinds)

		transform := Transform{}
		ok := ecs.Read(engine, id, &transform)
		if !ok {
			return
		}

		if window.Pressed(keybinds.Left) {
			transform.X -= 2.0
		}
		if window.Pressed(keybinds.Right) {
			transform.X += 2.0
		}
		if window.Pressed(keybinds.Up) {
			transform.Y += 2.0
		}
		if window.Pressed(keybinds.Down) {
			transform.Y -= 2.0
		}

		ecs.Write(engine, id, transform)
	})
}
