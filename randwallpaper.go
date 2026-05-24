package randwallpaper

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand/v2"
	"os"
	"path/filepath"
	"time"
)

type Option func(*config)

type config struct {
	count int
	seed  int64
}

func WithCount(n int) Option {
	return func(c *config) {
		c.count = n
	}
}

func WithSeed(seed int64) Option {
	return func(c *config) {
		c.seed = seed
	}
}

func generateImage(width, height int, rng *rand.Rand) *image.NRGBA {
	m := newMask(width, height, rng)
	p := newPath(m, rng)
	cm := generateColourmap(rng)
	path := p.Path()

	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	pathLen := len(path)
	for i, pt := range path {
		t := float64(i) / float64(pathLen-1)
		c := cm.mapAt(t)
		img.SetNRGBA(pt.X, pt.Y, color.NRGBA{R: c.R, G: c.G, B: c.B, A: 255})
	}
	return img
}

func Generate(width, height int, output string, opts ...Option) error {
	if width <= 0 || height <= 0 {
		return fmt.Errorf("randwallpaper: width and height must be positive, got %dx%d", width, height)
	}

	cfg := &config{
		count: 1,
		seed:  time.Now().UnixNano(),
	}
	for _, o := range opts {
		o(cfg)
	}
	if cfg.count <= 0 {
		return fmt.Errorf("randwallpaper: count must be positive, got %d", cfg.count)
	}

	saveSingle := func(idx int, rng *rand.Rand) error {
		var filename string
		if cfg.count == 1 {
			filename = output + ".png"
		} else {
			dir := output
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("randwallpaper: cannot create output directory: %w", err)
			}
			filename = filepath.Join(dir, fmt.Sprintf("%s_%d.png", filepath.Base(output), idx))
		}
		f, err := os.Create(filename)
		if err != nil {
			return fmt.Errorf("randwallpaper: cannot create file: %w", err)
		}
		defer f.Close()

		img := generateImage(width, height, rng)
		if err := png.Encode(f, img); err != nil {
			return fmt.Errorf("randwallpaper: png encode failed: %w", err)
		}
		return nil
	}

	for i := 1; i <= cfg.count; i++ {
		rng := rand.New(rand.NewPCG(uint64(cfg.seed), uint64(i)))
		if err := saveSingle(i, rng); err != nil {
			return err
		}
	}
	return nil
}
