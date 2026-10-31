package randwallpaper

import (
	"image"
	"math"
	"math/rand/v2"
	"testing"
)

func assertPathCoverage(t *testing.T, p path, width, height int) {
	t.Helper()

	points := p.Path()
	want := width * height
	if len(points) != want {
		t.Fatalf("path length = %d, want %d for %dx%d", len(points), want, width, height)
	}

	seen := make(map[image.Point]struct{}, want)
	for _, pt := range points {
		if pt.X < 0 || pt.X >= width || pt.Y < 0 || pt.Y >= height {
			t.Fatalf("path point %v is outside %dx%d bounds", pt, width, height)
		}
		if _, ok := seen[pt]; ok {
			t.Fatalf("duplicate path point %v for %dx%d", pt, width, height)
		}
		seen[pt] = struct{}{}
	}
	if len(seen) != want {
		t.Fatalf("path covers %d of %d pixels for %dx%d", len(seen), want, width, height)
	}
}

func pathForKind(kind uint8, m mask, seed uint64) path {
	rng := rand.New(rand.NewPCG(seed, 2))
	if kind%2 == 0 {
		return newEPWTPath(m, rng)
	}
	return newProbabilisticPath(m, rng)
}

func TestPathCoverageEdgeDimensions(t *testing.T) {
	dimensions := []image.Point{
		{1, 1},
		{1, 5},
		{5, 1},
		{1, 32},
		{32, 1},
		{2, 2},
		{3, 5},
		{5, 3},
		{8, 10},
		{10, 8},
		{13, 4},
		{4, 13},
	}

	for _, dim := range dimensions {
		for kind := uint8(0); kind < 2; kind++ {
			for seed := uint64(0); seed < 3; seed++ {
				maskRNG := rand.New(rand.NewPCG(seed+1, 1))
				m := newNormalMask(dim.X, dim.Y, maskRNG)
				assertPathCoverage(t, pathForKind(kind, m, seed+11), dim.X, dim.Y)
			}
		}
	}
}

func FuzzPathCoverage(f *testing.F) {
	f.Add(uint8(0), uint8(0), uint64(1), uint64(2), uint8(0))
	f.Add(uint8(0), uint8(9), uint64(3), uint64(4), uint8(1))
	f.Add(uint8(9), uint8(0), uint64(5), uint64(6), uint8(0))
	f.Add(uint8(6), uint8(6), uint64(7), uint64(8), uint8(1))
	f.Add(uint8(30), uint8(17), uint64(9), uint64(10), uint8(0))

	f.Fuzz(func(t *testing.T, widthRaw, heightRaw uint8, maskSeed, pathSeed uint64, kind uint8) {
		width := int(widthRaw%64) + 1
		height := int(heightRaw%64) + 1

		maskRNG := rand.New(rand.NewPCG(maskSeed, 1))
		m := newNormalMask(width, height, maskRNG)
		assertPathCoverage(t, pathForKind(kind, m, pathSeed), width, height)
	})
}

func FuzzColourmap(f *testing.F) {
	f.Add(uint64(1), uint32(0))
	f.Add(uint64(42), uint32(math.MaxUint32))
	f.Add(uint64(99), uint32(math.MaxUint32/2))

	f.Fuzz(func(t *testing.T, seed uint64, positionRaw uint32) {
		cm := generateColourmap(rand.New(rand.NewPCG(seed, 3)))
		position := float64(positionRaw) / float64(math.MaxUint32)

		c := cm.mapAt(position)
		if c.A != 255 {
			t.Fatalf("mapAt(%f) alpha = %d, want 255", position, c.A)
		}

		if got, want := cm.mapAt(0), cm.anchorColor(cm.anchors[0]); got != want {
			t.Fatalf("mapAt(0) = %v, want first anchor %v", got, want)
		}
		if got, want := cm.mapAt(1), cm.anchorColor(cm.anchors[len(cm.anchors)-1]); got != want {
			t.Fatalf("mapAt(1) = %v, want last anchor %v", got, want)
		}
	})
}

func TestColourmapDegenerateAnchors(t *testing.T) {
	maps := []colourmap{
		{anchors: []anchor{
			{pos: 0, r: 0.1, g: 0.2, b: 0.3},
		}},
		{anchors: []anchor{
			{pos: 0, r: 0.1, g: 0.2, b: 0.3},
			{pos: 1, r: 0.4, g: 0.5, b: 0.6},
		}},
		{anchors: []anchor{
			{pos: 0, r: 0.1, g: 0.2, b: 0.3},
			{pos: 0, r: 0.8, g: 0.1, b: 0.9},
			{pos: 0.5, r: 0.2, g: 0.6, b: 0.4},
			{pos: 0.5, r: 0.3, g: 0.7, b: 0.5},
			{pos: 1, r: 0.4, g: 0.5, b: 0.6},
		}},
	}

	for _, cm := range maps {
		for _, position := range []float64{math.Inf(-1), -1, 0, 0.25, 0.5, 0.75, 1, 2, math.Inf(1), math.NaN()} {
			c := cm.mapAt(position)
			if c.A != 255 {
				t.Fatalf("mapAt(%f) alpha = %d, want 255", position, c.A)
			}
		}
	}
}
