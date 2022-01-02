package main

//go:generate packer --input images --stats

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gommo/engine/asset"
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

	manSprite, err := spritesheet.Get(purpleGemPng)
	if err != nil {
		return
	}
	purpleGem := window.Bounds().Center()

	for !window.JustPressed(pixelgl.KeyEscape) {
		window.Clear(pixel.RGB(0, 0, 0))

		if window.Pressed(pixelgl.KeyLeft) {
			purpleGem.X -= 2.0
		}

		if window.Pressed(pixelgl.KeyRight) {
			purpleGem.X += 2.0
		}

		if window.Pressed(pixelgl.KeyUp) {
			purpleGem.Y += 2.0
		}

		if window.Pressed(pixelgl.KeyDown) {
			purpleGem.Y -= 2.0
		}

		manSprite.Draw(window, pixel.IM.Scaled(pixel.ZV, 0.1).Moved(purpleGem))

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
