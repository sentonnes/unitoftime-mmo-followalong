package proceduralgeneration

import (
	"github.com/ojrac/opensimplex-go"
	"math"
)

type Octave struct {
	Frequency float64
	Scale     float64
}

type NoiseMap struct {
	seed     int64
	noise    opensimplex.Noise
	octaves  []Octave
	exponent float64
}

func NewNoiseMap(seed int64, octaves []Octave, exponent float64) *NoiseMap {
	return &NoiseMap{
		seed:     seed,
		noise:    opensimplex.NewNormalized(seed),
		octaves:  octaves,
		exponent: exponent,
	}
}

func (n *NoiseMap) Get(x int, y int) float64 {
	ret := 0.0
	for i := range n.octaves {
		xNoise := n.octaves[i].Frequency * float64(x)
		yNoise := n.octaves[i].Frequency * float64(y)
		ret += n.octaves[i].Scale * n.noise.Eval2(xNoise, yNoise)
	}

	ret = math.Pow(ret, n.exponent)
	return ret
}
