package main

import "github.com/faiface/pixel/pixelgl"

type Keybinds struct {
	Up, Down, Left, Right pixelgl.Button
}

func (keybinds *Keybinds) ComponentSet(val interface{}) {
	*keybinds = val.(Keybinds)
}

var AWSDKeybinds = Keybinds{Up: pixelgl.KeyW, Down: pixelgl.KeyS, Left: pixelgl.KeyA, Right: pixelgl.KeyD}
var ArrowKeybinds = Keybinds{Up: pixelgl.KeyUp, Down: pixelgl.KeyDown, Left: pixelgl.KeyLeft, Right: pixelgl.KeyRight}
