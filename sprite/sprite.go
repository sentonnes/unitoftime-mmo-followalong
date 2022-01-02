package sprite

import (
	"github.com/faiface/pixel"
	"image"
	_ "image/png"
	"os"
	"path"
)

const meatPng = "meat.png"
const pngDirectory = "./pngs"

func MeatSprite() (*pixel.Sprite, error) {
	sprite, err := getSprite(meatPng)
	return sprite, err
}

func getSprite(pngName string) (*pixel.Sprite, error) {
	pathToPng := path.Join(pngDirectory, pngName)
	file, err := os.Open(pathToPng)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	pic := pixel.PictureDataFromImage(img)

	return pixel.NewSprite(pic, pic.Bounds()), nil
}
