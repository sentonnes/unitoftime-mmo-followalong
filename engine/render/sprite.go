package render

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"gommo/engine/ecs"
	"gommo/engine/physics"
)

type Sprite struct {
	*pixel.Sprite
}

func (sprite *Sprite) ComponentSet(val interface{}) { *sprite = val.(Sprite) }

type Keybinds struct {
	Up, Down, Left, Right pixelgl.Button
}

var AWSDKeybinds = Keybinds{Up: pixelgl.KeyW, Down: pixelgl.KeyS, Left: pixelgl.KeyA, Right: pixelgl.KeyD}
var ArrowKeybinds = Keybinds{Up: pixelgl.KeyUp, Down: pixelgl.KeyDown, Left: pixelgl.KeyLeft, Right: pixelgl.KeyRight}

func (t *Keybinds) ComponentSet(val interface{}) { *t = val.(Keybinds) }

func DrawSprites(win *pixelgl.Window, engine *ecs.Engine) {
	ecs.Each(engine, Sprite{}, func(id ecs.Id, a interface{}) {
		sprite := a.(Sprite)

		transform := physics.Transform{}
		ok := ecs.Read(engine, id, &transform)
		if !ok {
			return
		}

		pos := pixel.V(transform.X, transform.Y)
		sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, 2.0).Moved(pos))
	})
}

func CaptureInput(win *pixelgl.Window, engine *ecs.Engine) {
	ecs.Each(engine, Keybinds{}, func(id ecs.Id, a interface{}) {
		keybinds := a.(Keybinds)

		input := physics.Input{}
		ok := ecs.Read(engine, id, &input)
		if !ok {
			return
		}

		input.Left = false
		input.Right = false
		input.Up = false
		input.Down = false

		if win.Pressed(keybinds.Left) {
			input.Left = true
		}
		if win.Pressed(keybinds.Right) {
			input.Right = true
		}
		if win.Pressed(keybinds.Up) {
			input.Up = true
		}
		if win.Pressed(keybinds.Down) {
			input.Down = true
		}

		ecs.Write(engine, id, input)
	})
}
