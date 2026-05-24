package randwallpaper

import (
	"fmt"
	"image"
	"image/png"
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
)

func TestSaltPepperMask(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	m := newSaltPepperMask(50, 30, rng)
	if m.Width() != 50 || m.Height() != 30 {
		t.Fatalf("expected 50x30, got %dx%d", m.Width(), m.Height())
	}
	for y := range 30 {
		for x := range 50 {
			v := m.At(x, y)
			if v != 0 && v != 1 {
				t.Fatalf("saltPepperMask value at (%d,%d) = %f, expected 0 or 1", x, y, v)
			}
		}
	}
}

func TestSaltPepperMaskDeterministic(t *testing.T) {
	rng1 := rand.New(rand.NewPCG(42, 1))
	m1 := newSaltPepperMask(20, 20, rng1)
	rng2 := rand.New(rand.NewPCG(42, 1))
	m2 := newSaltPepperMask(20, 20, rng2)
	for y := range 20 {
		for x := range 20 {
			if m1.At(x, y) != m2.At(x, y) {
				t.Fatalf("deterministic mask mismatch at (%d,%d)", x, y)
			}
		}
	}
}

func TestNormalMask(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	m := newNormalMask(50, 30, rng)
	if m.Width() != 50 || m.Height() != 30 {
		t.Fatalf("expected 50x30, got %dx%d", m.Width(), m.Height())
	}
	for y := range 30 {
		for x := range 50 {
			v := m.At(x, y)
			if v < 0 {
				t.Fatalf("normalMask value at (%d,%d) = %f, expected >= 0", x, y, v)
			}
		}
	}
}

func TestNormalMaskDeterministic(t *testing.T) {
	rng1 := rand.New(rand.NewPCG(99, 2))
	m1 := newNormalMask(20, 20, rng1)
	rng2 := rand.New(rand.NewPCG(99, 2))
	m2 := newNormalMask(20, 20, rng2)
	for y := range 20 {
		for x := range 20 {
			if m1.At(x, y) != m2.At(x, y) {
				t.Fatalf("deterministic mask mismatch at (%d,%d)", x, y)
			}
		}
	}
}

func TestGaussianBlobMask(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	m := newGaussianBlobMask(50, 30, rng)
	if m.Width() != 50 || m.Height() != 30 {
		t.Fatalf("expected 50x30, got %dx%d", m.Width(), m.Height())
	}
	hasNonZero := false
	for y := range 30 {
		for x := range 50 {
			v := m.At(x, y)
			if v < 0 || v > 1 {
				t.Fatalf("gaussianBlobMask value at (%d,%d) = %f, expected [0,1]", x, y, v)
			}
			if v > 0 {
				hasNonZero = true
			}
		}
	}
	if !hasNonZero {
		t.Fatal("gaussianBlobMask should have some non-zero values")
	}
}

func TestGaussianBlobMaskDeterministic(t *testing.T) {
	rng1 := rand.New(rand.NewPCG(77, 3))
	m1 := newGaussianBlobMask(20, 20, rng1)
	rng2 := rand.New(rand.NewPCG(77, 3))
	m2 := newGaussianBlobMask(20, 20, rng2)
	for y := range 20 {
		for x := range 20 {
			if m1.At(x, y) != m2.At(x, y) {
				t.Fatalf("deterministic gaussianBlobMask mismatch at (%d,%d)", x, y)
			}
		}
	}
}

func TestEPWTPathCoversAllPixels(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	m := newNormalMask(10, 10, rng)
	p := newEPWTPath(m, rand.New(rand.NewPCG(99, 1)))
	path := p.Path()
	expected := 10 * 10
	if len(path) != expected {
		t.Fatalf("EPWT path length = %d, expected %d", len(path), expected)
	}
	seen := make(map[image.Point]bool)
	for _, pt := range path {
		if seen[pt] {
			t.Fatalf("duplicate point (%d,%d) in EPWT path", pt.X, pt.Y)
		}
		seen[pt] = true
	}
}

func TestProbabilisticPathCoversAllPixels(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	m := newNormalMask(10, 10, rng)
	p := newProbabilisticPath(m, rand.New(rand.NewPCG(99, 1)))
	path := p.Path()
	expected := 10 * 10
	if len(path) != expected {
		t.Fatalf("Probabilistic path length = %d, expected %d", len(path), expected)
	}
	seen := make(map[image.Point]bool)
	for _, pt := range path {
		if seen[pt] {
			t.Fatalf("duplicate point (%d,%d) in Probabilistic path", pt.X, pt.Y)
		}
		seen[pt] = true
	}
}

func TestEPWTPathAllPixelsVisitedExactly(t *testing.T) {
	for seed := range uint64(5) {
		rng := rand.New(rand.NewPCG(42, uint64(seed)))
		m := newNormalMask(8, 8, rng)
		p := newEPWTPath(m, rand.New(rand.NewPCG(99, uint64(seed))))
		path := p.Path()
		if len(path) != 64 {
			t.Fatalf("seed %d: path length %d, expected 64", seed, len(path))
		}
		seen := make(map[image.Point]bool)
		for _, pt := range path {
			seen[pt] = true
		}
		if len(seen) != 64 {
			t.Fatalf("seed %d: only %d unique points visited", seed, len(seen))
		}
	}
}

func TestProbabilisticPathAllPixelsVisitedExactly(t *testing.T) {
	for seed := range uint64(5) {
		rng := rand.New(rand.NewPCG(42, uint64(seed)))
		m := newNormalMask(8, 8, rng)
		p := newProbabilisticPath(m, rand.New(rand.NewPCG(99, uint64(seed))))
		path := p.Path()
		if len(path) != 64 {
			t.Fatalf("seed %d: path length %d, expected 64", seed, len(path))
		}
		seen := make(map[image.Point]bool)
		for _, pt := range path {
			seen[pt] = true
		}
		if len(seen) != 64 {
			t.Fatalf("seed %d: only %d unique points visited", seed, len(seen))
		}
	}
}

func TestPathDeterministic(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	m := newNormalMask(10, 10, rng)

	p1 := newEPWTPath(m, rand.New(rand.NewPCG(77, 1)))
	path1 := p1.Path()

	p2 := newEPWTPath(m, rand.New(rand.NewPCG(77, 1)))
	path2 := p2.Path()

	if len(path1) != len(path2) {
		t.Fatalf("path lengths differ: %d vs %d", len(path1), len(path2))
	}
	for i := range path1 {
		if path1[i] != path2[i] {
			t.Fatalf("path differs at index %d: (%d,%d) vs (%d,%d)",
				i, path1[i].X, path1[i].Y, path2[i].X, path2[i].Y)
		}
	}
}

func TestColourmap(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	cm := generateColourmap(rng)
	if len(cm.anchors) < 3 || len(cm.anchors) > 7 {
		t.Fatalf("expected 3-7 anchor points, got %d", len(cm.anchors))
	}
	if cm.anchors[0].pos != 0 {
		t.Fatalf("first anchor pos = %f, expected 0", cm.anchors[0].pos)
	}
	last := cm.anchors[len(cm.anchors)-1].pos
	if last != 1 {
		t.Fatalf("last anchor pos = %f, expected 1", last)
	}
	for i, a := range cm.anchors {
		if a.r < 0 || a.r > 1 || a.g < 0 || a.g > 1 || a.b < 0 || a.b > 1 {
			t.Fatalf("anchor %d has out-of-range colour: (%f,%f,%f)", i, a.r, a.g, a.b)
		}
	}
}

func TestColourmapOutput(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	cm := generateColourmap(rng)
	for _, tc := range []float64{0, 0.25, 0.5, 0.75, 1} {
		c := cm.mapAt(tc)
		if c.A != 255 {
			t.Fatalf("mapAt(%f) alpha = %d, expected 255", tc, c.A)
		}
	}
	c0 := cm.mapAt(0)
	c1 := cm.mapAt(1)
	if c0.R != uint8(cm.anchors[0].r*255) ||
		c0.G != uint8(cm.anchors[0].g*255) ||
		c0.B != uint8(cm.anchors[0].b*255) {
		t.Fatalf("mapAt(0) doesn't match first anchor")
	}
	last := cm.anchors[len(cm.anchors)-1]
	if c1.R != uint8(last.r*255) ||
		c1.G != uint8(last.g*255) ||
		c1.B != uint8(last.b*255) {
		t.Fatalf("mapAt(1) doesn't match last anchor")
	}
}

func TestClampF64(t *testing.T) {
	tests := []struct {
		in, out float64
	}{
		{-0.5, 0},
		{0, 0},
		{0.5, 0.5},
		{1, 1},
		{1.5, 1},
	}
	for _, tc := range tests {
		got := clampF64(tc.in)
		if got != tc.out {
			t.Errorf("clampF64(%f) = %f, expected %f", tc.in, got, tc.out)
		}
	}
}

func TestGeneratePNG(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "testimg")
	err := Generate(50, 30, out, WithSeed(42))
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}
	f, err := os.Open(out + ".png")
	if err != nil {
		t.Fatalf("cannot open generated file: %v", err)
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("decoding PNG failed: %v", err)
	}
	b := img.Bounds()
	if b.Dx() != 50 || b.Dy() != 30 {
		t.Fatalf("image bounds = %dx%d, expected 50x30", b.Dx(), b.Dy())
	}
}

func TestGenerateCount(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "batch")
	err := Generate(20, 20, out, WithCount(3), WithSeed(42))
	if err != nil {
		t.Fatalf("Generate batch failed: %v", err)
	}
	for i := 1; i <= 3; i++ {
		filename := filepath.Join(dir, "batch", fmt.Sprintf("batch_%d.png", i))
		f, err := os.Open(filename)
		if err != nil {
			t.Fatalf("missing batch file %d: %v", i, err)
		}
		_, err = png.Decode(f)
		f.Close()
		if err != nil {
			t.Fatalf("batch file %d is not valid PNG: %v", i, err)
		}
	}
}

func TestGenerateSeedReproducibility(t *testing.T) {
	dir1 := t.TempDir()
	dir2 := t.TempDir()

	err := Generate(30, 30, filepath.Join(dir1, "img"), WithSeed(123))
	if err != nil {
		t.Fatal(err)
	}
	err = Generate(30, 30, filepath.Join(dir2, "img"), WithSeed(123))
	if err != nil {
		t.Fatal(err)
	}

	f1, _ := os.Open(filepath.Join(dir1, "img.png"))
	defer f1.Close()
	f2, _ := os.Open(filepath.Join(dir2, "img.png"))
	defer f2.Close()

	img1, _ := png.Decode(f1)
	img2, _ := png.Decode(f2)

	b := img1.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)
			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				t.Fatalf("pixel (%d,%d) differs between runs: (%d,%d,%d,%d) vs (%d,%d,%d,%d)",
					x, y, r1, g1, b1, a1, r2, g2, b2, a2)
			}
		}
	}
}

func TestGenerateInvalidDimensions(t *testing.T) {
	tests := []struct {
		w, h int
	}{
		{0, 100},
		{100, 0},
		{-1, 100},
		{100, -1},
	}
	for _, tc := range tests {
		err := Generate(tc.w, tc.h, t.TempDir()+"/out", WithSeed(42))
		if err == nil {
			t.Fatalf("expected error for %dx%d, got nil", tc.w, tc.h)
		}
	}
}

func TestGenerateNegativeCount(t *testing.T) {
	err := Generate(10, 10, t.TempDir()+"/out", WithCount(0))
	if err == nil {
		t.Fatal("expected error for count=0, got nil")
	}
	err = Generate(10, 10, t.TempDir()+"/out", WithCount(-1))
	if err == nil {
		t.Fatal("expected error for count=-1, got nil")
	}
}

func TestGenerate1x1(t *testing.T) {
	dir := t.TempDir()
	err := Generate(1, 1, filepath.Join(dir, "one"), WithSeed(42))
	if err != nil {
		t.Fatalf("Generate 1x1 failed: %v", err)
	}
	f, _ := os.Open(filepath.Join(dir, "one.png"))
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		t.Fatalf("decoding 1x1 PNG failed: %v", err)
	}
	b := img.Bounds()
	if b.Dx() != 1 || b.Dy() != 1 {
		t.Fatalf("1x1 image bounds = %dx%d", b.Dx(), b.Dy())
	}
}

func TestColourmapSmoothInterpolation(t *testing.T) {
	rng := rand.New(rand.NewPCG(42, 1))
	cm := generateColourmap(rng)

	steps := 100
	prev := cm.mapAt(0)
	for i := 1; i <= steps; i++ {
		pos := float64(i) / float64(steps)
		cur := cm.mapAt(pos)
		dr := int(cur.R) - int(prev.R)
		dg := int(cur.G) - int(prev.G)
		db := int(cur.B) - int(prev.B)
		dist := math.Sqrt(float64(dr*dr + dg*dg + db*db))
		maxDist := 255.0 * math.Sqrt(3)
		if dist > maxDist {
			t.Fatalf("suspicious jump at pos=%f: distance=%f", pos, dist)
		}
		prev = cur
	}
}

func TestNewSaltPepperMask(t *testing.T) {
	rng := rand.New(rand.NewPCG(1, 1))
	m := newSaltPepperMask(100, 100, rng)
	if m.Width() != 100 || m.Height() != 100 {
		t.Fatalf("wrong dimensions")
	}
	count := 0
	for y := range 100 {
		for x := range 100 {
			if m.At(x, y) == 1 {
				count++
			}
		}
	}
	if count == 0 || count == 10000 {
		t.Fatal("saltPepperMask should have a mix of 0s and 1s")
	}
}

func BenchmarkEPWTPath(b *testing.B) {
	for range b.N {
		rng := rand.New(rand.NewPCG(42, 1))
		m := newNormalMask(64, 64, rng)
		p := newEPWTPath(m, rand.New(rand.NewPCG(99, 1)))
		_ = p.Path()
	}
}

func BenchmarkProbabilisticPath(b *testing.B) {
	for range b.N {
		rng := rand.New(rand.NewPCG(42, 1))
		m := newNormalMask(64, 64, rng)
		p := newProbabilisticPath(m, rand.New(rand.NewPCG(99, 1)))
		_ = p.Path()
	}
}

func BenchmarkGenerate(b *testing.B) {
	for range b.N {
		_ = generateImage(128, 128, rand.New(rand.NewPCG(42, 1)))
	}
}
