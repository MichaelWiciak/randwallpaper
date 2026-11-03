package randwallpaper

import (
	"math"
	"math/rand/v2"
)

type mask interface {
	At(x, y int) float64
	Width() int
	Height() int
}

type saltPepperMask struct {
	data   []float64
	width  int
	height int
}

func newSaltPepperMask(width, height int, rng *rand.Rand) *saltPepperMask {
	data := make([]float64, width*height)
	for i := range data {
		data[i] = float64(rng.IntN(2))
	}
	return &saltPepperMask{data: data, width: width, height: height}
}

func (m *saltPepperMask) At(x, y int) float64 { return m.data[y*m.width+x] }
func (m *saltPepperMask) Width() int          { return m.width }
func (m *saltPepperMask) Height() int         { return m.height }

type normalMask struct {
	data   []float64
	width  int
	height int
}

func newNormalMask(width, height int, rng *rand.Rand) *normalMask {
	data := make([]float64, width*height)
	minVal := 0.0
	for i := range data {
		v := rng.NormFloat64()
		data[i] = v
		if v < minVal {
			minVal = v
		}
	}
	if minVal < 0 {
		for i := range data {
			data[i] -= minVal
		}
	}
	return &normalMask{data: data, width: width, height: height}
}

func (m *normalMask) At(x, y int) float64 { return m.data[y*m.width+x] }
func (m *normalMask) Width() int          { return m.width }
func (m *normalMask) Height() int         { return m.height }

type gaussianBlobMask struct {
	data   []float64
	width  int
	height int
}

func newGaussianBlobMask(width, height int, rng *rand.Rand) *gaussianBlobMask {
	n := width * height
	sqrtN := math.Sqrt(float64(n))
	centerMax := max(1, int(0.5*sqrtN))
	sigmaMax := max(1, int(0.2*sqrtN))
	nCenters := rng.IntN(centerMax) + 1
	sigma := rng.IntN(sigmaMax) + 1

	centers := make([]uint8, n)
	for range nCenters {
		cx := rng.IntN(width)
		cy := rng.IntN(height)
		centers[cy*width+cx] = 255
	}

	blurred := gaussianBlur(centers, width, height, float64(sigma))
	data := make([]float64, n)
	maxVal := 0.0
	for i, v := range blurred {
		f := float64(v)
		data[i] = f
		if f > maxVal {
			maxVal = f
		}
	}
	if maxVal > 0 {
		for i := range data {
			data[i] /= maxVal
		}
	}
	return &gaussianBlobMask{data: data, width: width, height: height}
}

func (m *gaussianBlobMask) At(x, y int) float64 { return m.data[y*m.width+x] }
func (m *gaussianBlobMask) Width() int          { return m.width }
func (m *gaussianBlobMask) Height() int         { return m.height }

func newMask(width, height int, rng *rand.Rand) mask {
	switch rng.IntN(3) {
	case 0:
		return newSaltPepperMask(width, height, rng)
	case 1:
		return newNormalMask(width, height, rng)
	default:
		return newGaussianBlobMask(width, height, rng)
	}
}
