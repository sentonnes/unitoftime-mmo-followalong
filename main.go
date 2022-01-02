package main

//go:generate packer --input images --stats

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gommo/engine/asset"
	"gommo/engine/render"
	_ "image/png"
	"os"
)

func main() {
	pixelgl.Run(runGame)
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
)

var load *asset.Load
var window *pixelgl.Window

func runGame() {
	setupGame()
	runGameLoop()
}

func runGameLoop() {
	spritesheet, err := load.Spritesheet(packedJson)

	purpleGemSprite, err := spritesheet.Get(purpleGemPng)
	if err != nil {
		return
	}
	purpleGemPosition := window.Bounds().Center()

	redGemSprite, err := spritesheet.Get(redGemPng)
	if err != nil {
		return
	}
	redGemPosition := window.Bounds().Center()

	people := make([]Person, 0)
	newPerson := NewPerson(purpleGemSprite, purpleGemPosition, Keybinds{Up: pixelgl.KeyUp, Down: pixelgl.KeyDown, Left: pixelgl.KeyLeft, Right: pixelgl.KeyRight})
	people = append(people, newPerson)
	newPerson = NewPerson(redGemSprite, redGemPosition, Keybinds{Up: pixelgl.KeyW, Down: pixelgl.KeyS, Left: pixelgl.KeyA, Right: pixelgl.KeyD})
	people = append(people, newPerson)

	camera := render.NewCamera(window, 0, 0)
	zoomSpeed := 0.1

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
		for i := range people {
			people[i].Draw(window)
		}
		window.SetMatrix(pixel.IM)

		window.Update()
	}
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
	if err != nil {
		panic(err)
	}

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
