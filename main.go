package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gommo/sprite"
	_ "image/png"
)

func main() {
	pixelgl.Run(runGame)
}

const windowsMaxX = 1024
const windowsMaxY = 768
const windowsName = "MMO"
const windowsVSync = true
const windowsResizable = true

func runGame() {
	window := setupGame()
	runGameLoop(window)
}

func runGameLoop(win *pixelgl.Window) {
	manSprite, err := sprite.MeatSprite()
	if err != nil {
		panic(err)
	}
	manPosition := win.Bounds().Center()

	for !win.JustPressed(pixelgl.KeyEscape) {
		win.Clear(pixel.RGB(0, 0, 0))

		if win.Pressed(pixelgl.KeyLeft) {
			manPosition.X -= 2.0
		}

		if win.Pressed(pixelgl.KeyRight) {
			manPosition.X += 2.0
		}

		if win.Pressed(pixelgl.KeyUp) {
			manPosition.Y += 2.0
		}

		if win.Pressed(pixelgl.KeyDown) {
			manPosition.Y -= 2.0
		}

		manSprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(manPosition))

		win.Update()
	}
}

func setupGame() *pixelgl.Window {
	cfg := getWindowsConfig()

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	win.SetSmooth(false)
	return win
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
