package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type Person struct {
	Sprite        *pixel.Sprite
	Position      pixel.Vec
	Keybinds      Keybinds
	MovementSpeed float64
}

func NewPerson(sprite *pixel.Sprite, position pixel.Vec, keybinds Keybinds, movementSpeed float64) Person {
	return Person{sprite, position, keybinds, movementSpeed}
}

func (person *Person) Draw(window *pixelgl.Window) {
	person.Sprite.Draw(window, pixel.IM.Scaled(pixel.ZV, 1).Moved(person.Position))
}

func (person *Person) HandleInput(window *pixelgl.Window) {
	if window.Pressed(person.Keybinds.Left) {
		person.Position.X -= person.MovementSpeed
	}
	if window.Pressed(person.Keybinds.Right) {
		person.Position.X += person.MovementSpeed
	}
	if window.Pressed(person.Keybinds.Up) {
		person.Position.Y += person.MovementSpeed
	}
	if window.Pressed(person.Keybinds.Down) {
		person.Position.Y -= person.MovementSpeed
	}
}
