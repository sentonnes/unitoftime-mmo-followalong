package asset

import (
	"encoding/json"
	"github.com/faiface/pixel"
	"github.com/unitoftime/packer"
	"image"
	_ "image/png"
	"io/fs"
	"io/ioutil"
)

type Load struct {
	filesystem fs.FS
}

func NewLoad(filesystem fs.FS) *Load {
	return &Load{filesystem: filesystem}
}

func (load *Load) Open(path string) (fs.File, error) {
	return load.filesystem.Open(path)
}

func (load *Load) Image(path string) (image.Image, error) {
	file, err := load.filesystem.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, err
}

func (load *Load) Sprite(path string) (*pixel.Sprite, error) {
	img, err := load.Image(path)
	if err != nil {
		return nil, err
	}

	pic := pixel.PictureDataFromImage(img)

	return pixel.NewSprite(pic, pic.Bounds()), nil
}

func (load *Load) Json(path string, data interface{}) error {
	file, err := load.filesystem.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, &data)
}

func (load *Load) Spritesheet(path string) (*Spritesheet, error) {
	//load the json
	serializedSpritesheet := packer.SerializedSpritesheet{}
	err := load.Json(path, &serializedSpritesheet)
	if err != nil {
		return nil, err
	}

	// load the image
	img, err := load.Image(serializedSpritesheet.ImageName)
	if err != nil {
		return nil, err
	}
	pic := pixel.PictureDataFromImage(img)

	// create the spritesheet
	bounds := pic.Bounds()
	lookup := make(map[string]*pixel.Sprite)
	for k, v := range serializedSpritesheet.Frames {
		rect := pixel.R(
			v.Frame.X,
			bounds.H()-v.Frame.Y,
			v.Frame.X+v.Frame.W,
			bounds.H()-(v.Frame.Y+v.Frame.H)).Norm()

		lookup[k] = pixel.NewSprite(pic, rect)
	}

	return NewSpritesheet(pic, lookup), nil
}
