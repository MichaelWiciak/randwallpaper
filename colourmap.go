package randwallpaper

import (
	"image/color"
	"math/rand/v2"
	"sort"

	"google.golang.org/grpc/credentials/local"
)

type anchor struct {
	pos    float64
	r, g, b float64
}

type colourmap struct {
	anchors []anchor
}

func generateColourmap(rng *rand.Rand) colourmap {
	n := rng.IntN(5) + 3
	anchors := make([]anchor, n)
	for i := range anchors {
		anchors[i] = anchor{
			pos: rng.Float64(),
			r:   rng.Float64(),
			g:   rng.Float64(),
			b:   rng.Float64(),
		}
	}
	anchors[0].pos = 0
	anchors[n-1].pos = 1
	sort.Slice(anchors, func(i, j int) bool { return anchors[i].pos < anchors[j].pos })
	return colourmap{anchors: anchors}
}

func (cm colourmap) mapAt(t float64) color.RGBA {
	if t <= 0 {
		return cm.anchorColor(cm.anchors[0])
	}
	if t >= 1 {
		return cm.anchorColor(cm.anchors[len(cm.anchors)-1])
	}
	for i := 0; i < len(cm.anchors)-1; i++ {
		a, b := cm.anchors[i], cm.anchors[i+1]
		if t >= a.pos && t <= b.pos {
			den := b.pos - a.pos
			if den <= 0 {
				return cm.anchorColor(a)
			}
			local := (t - a.pos) / den
			local = local * local * (3 - 2*local)
			return color.RGBA{
				R: uint8(clampF64(a.r + (b.r-a.r)*local) * 255),
				G: uint8(clampF64(a.g + (b.g-a.g)*local) * 255),
				B: uint8(clampF64(a.b + (b.b-a.b)*local) * 255),
				A: 255,
			}
		}
	}
	return cm.anchorColor(cm.anchors[len(cm.anchors)-1])
}

func (cm colourmap) anchorColor(a anchor) color.RGBA {
	return color.RGBA{
		R: uint8(clampF64(a.r) * 255),
		G: uint8(clampF64(a.g) * 255),
		B: uint8(clampF64(a.b) * 255),
		A: 255,
	}
}

func clampF64(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
