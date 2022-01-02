package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"image"
	_ "image/png"
	"os"
	"path"
)

func main() {
	pixelgl.Run(runGame)
}

const windowsMaxX = 1024
const windowsMaxY = 768
const windowsName = "MMO"
const windowsVSync = true
const windowsResizable = true
const sprite = "meat.PNG"

func runGame() {
	window := setupGame()
	runGameLoop(window)
}

func runGameLoop(win *pixelgl.Window) {
	manSprite, err := getManSprite()
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

func getManSprite() (*pixel.Sprite, error) {
	sprite, err := getSprite(path.Join("./pngs", sprite))
	if err != nil {
		return nil, err
	}

	return sprite, nil
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

func getSprite(path string) (*pixel.Sprite, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	pic := pixel.PictureDataFromImage(img)

	return pixel.NewSprite(pic, pic.Bounds()), nil
}
