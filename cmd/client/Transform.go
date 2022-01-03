package main

type Transform struct {
	X float64
	Y float64
}

func (transform *Transform) ComponentSet(val interface{}) { *transform = val.(Transform) }
