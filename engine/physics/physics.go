package physics

import (
	"gommo/engine/ecs"
)

type Transform struct {
	X float64
	Y float64
}

func (transform *Transform) ComponentSet(val interface{}) { *transform = val.(Transform) }

type Input struct {
	Up, Down, Left, Right bool
}

func (input *Input) ComponentSet(val interface{}) {
	*input = val.(Input)
}

func HandleInput(engine *ecs.Engine) {
	ecs.Each(engine, Input{}, func(id ecs.Id, a interface{}) {
		input := a.(Input)

		transform := Transform{}
		ok := ecs.Read(engine, id, &transform)
		if !ok {
			return
		}

		if input.Left {
			transform.X -= 2.0
		}
		if input.Right {
			transform.X += 2.0
		}
		if input.Up {
			transform.Y += 2.0
		}
		if input.Down {
			transform.Y -= 2.0
		}

		ecs.Write(engine, id, transform)
	})
}
