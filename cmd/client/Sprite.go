package main

import "github.com/faiface/pixel"

type Sprite struct {
	*pixel.Sprite
}

func (sprite *Sprite) ComponentSet(val interface{}) { *sprite = val.(Sprite) }
