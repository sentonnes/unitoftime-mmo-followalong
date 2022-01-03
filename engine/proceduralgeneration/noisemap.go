package proceduralgeneration

import (
	"github.com/ojrac/opensimplex-go"
)

type NoiseMap struct {
	seed     int64
	noise    opensimplex.Noise
	exponent float64
}

func NewNoiseMap(seed int64, exponent float64) *NoiseMap {
	return &NoiseMap{
		seed:     seed,
		noise:    opensimplex.NewNormalized(seed),
		exponent: exponent,
	}
}

func (NoiseMap *NoiseMap) Get(x int, y int) float64 {
	frequency := 0.1
	yNoise := frequency * float64(y)
	xNoise := frequency * float64(x)
	ret := NoiseMap.noise.Eval2(xNoise, yNoise)
	return ret
}
