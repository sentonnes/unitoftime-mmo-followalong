package asset

import (
	"errors"
	"github.com/faiface/pixel"
)

type Spritesheet struct {
	picture pixel.Picture
	lookup  map[string]*pixel.Sprite
}

func NewSpritesheet(pic pixel.Picture, lookup map[string]*pixel.Sprite) *Spritesheet {
	return &Spritesheet{pic, lookup}
}

func (s *Spritesheet) Get(name string) (*pixel.Sprite, error) {
	sprite, ok := s.lookup[name]
	if !ok {
		return nil, errors.New("invalid sprite name")
	}
	return sprite, nil
}

func (s *Spritesheet) Picture() pixel.Picture {
	return s.picture
}
