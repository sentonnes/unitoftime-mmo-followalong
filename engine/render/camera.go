package render

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"math"
)

type Camera struct {
	window   *pixelgl.Window
	Position pixel.Vec
	Zoom     float64
	matrix   pixel.Matrix
}

func NewCamera(window *pixelgl.Window, x float64, y float64) *Camera {
	return &Camera{window, pixel.V(x, y), 1.0, pixel.IM}
}

func (camera *Camera) Update() {
	screenCenter := camera.window.Bounds().Center()

	movePosition := pixel.V(
		math.Floor(-camera.Position.X),
		math.Floor(-camera.Position.Y)).Add(screenCenter)

	camera.matrix = pixel.IM.Moved(movePosition).Scaled(screenCenter, camera.Zoom)
}

func (camera *Camera) Matrix() pixel.Matrix {
	return camera.matrix
}
