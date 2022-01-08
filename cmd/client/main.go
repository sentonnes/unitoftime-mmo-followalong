package main

//go:generate packer --input images --stats

import (
	"context"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	mmo "gommo"
	"gommo/engine/asset"
	"gommo/engine/ecs"
	"gommo/engine/physics"
	"gommo/engine/render"
	"gommo/engine/tilemap"
	_ "image/png"
	"log"
	"net"
	"nhooyr.io/websocket"
	"os"
	"time"
)

func main() {
	conn := createConnection()
	go sendCounterToServer(conn)
	pixelgl.Run(runGame)
}

func sendCounterToServer(conn net.Conn) {
	func() {
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
}

func createConnection() net.Conn {
	url := "ws://localhost:8000"
	ctx := context.Background()
	c, resp, err := websocket.Dial(ctx, url, nil)
	check(err)

	log.Println("Connection Response:", resp)

	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)
	return conn
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
	waterPng         = "water.png"
	sandPng          = "sand.png"
	grassPng         = "grass.png"
)

var load *asset.Load
var spritesheet *asset.Spritesheet
var window *pixelgl.Window
var engine *ecs.Engine
var playerId ecs.Id

func runGame() {
	setupGame()
	runGameLoop()
}

func runGameLoop() {
	tmap, purpleGemId, redGemId := mmo.LoadGame(engine)
	playerId = purpleGemId
	createPeople(spritesheet, purpleGemId, redGemId)
	tmapRenderer := createTileMapRender(tmap)
	gameLoop(tmapRenderer)
}

func gameLoop(tmapRender *render.TilemapRender) {
	camera, zoomSpeed := createCamera()
	quit := ecs.Signal{}
	quit.Set(false)

	inputSystems := []ecs.System{
		{"Clear", func(dt time.Duration) {
			window.Clear(pixel.RGB(0, 0, 0))

			scroll := window.MouseScroll()
			if scroll.Y != 0 {
				camera.Zoom += zoomSpeed * scroll.Y
			}

			if window.JustPressed(pixelgl.KeyEscape) {
				quit.Set(true)
			}
		}},
		{"CaptureInput", func(dt time.Duration) {
			render.CaptureInput(window, engine)
		}},
	}

	physicsSystems := mmo.CreatePhysicsSystems(engine)

	renderSystems := []ecs.System{
		{"UpdateCamera", func(dt time.Duration) {
			transform := physics.Transform{}
			ok := ecs.Read(engine, playerId, &transform)
			if ok {
				camera.Position = pixel.V(transform.X, transform.Y)
			}
			camera.Update()
		}},
		{"Draw", func(dt time.Duration) {
			window.SetMatrix(camera.Matrix())
			tmapRender.Draw(window)

			render.DrawSprites(window, engine)

			window.SetMatrix(pixel.IM)
		}},
		{"UpdateWindow", func(dt time.Duration) {
			window.Update()
		}},
	}

	ecs.RunGame(inputSystems, physicsSystems, renderSystems, &quit)
}

func createTileMapRender(tmap *tilemap.Tilemap) *render.TilemapRender {
	grassTile, err := spritesheet.Get(grassPng)
	check(err)
	sandTile, err := spritesheet.Get(sandPng)
	check(err)
	waterTile, err := spritesheet.Get(waterPng)
	check(err)

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

func createPeople(spritesheet *asset.Spritesheet, purpleGemId ecs.Id, redGemId ecs.Id) {
	purpleGemSprite, err := spritesheet.Get(purpleGemPng)
	check(err)
	ecs.Write(engine, purpleGemId, render.Sprite{Sprite: purpleGemSprite})
	ecs.Write(engine, purpleGemId, render.AWSDKeybinds)

	redGemSprite, err := spritesheet.Get(redGemPng)
	check(err)
	ecs.Write(engine, redGemId, render.Sprite{Sprite: redGemSprite})
	ecs.Write(engine, redGemId, render.ArrowKeybinds)
}

func setupGame() {
	setupAssets()
	setupEngine()
	setupWindow()
}

func setupEngine() {
	engine = ecs.NewEngine()
}

func setupAssets() {
	load = asset.NewLoad(os.DirFS("./"))
	var err error
	spritesheet, err = load.Spritesheet("packed.json")
	check(err)
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
