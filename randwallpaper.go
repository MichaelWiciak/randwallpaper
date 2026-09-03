// Package randwallpaper generates unique, procedural wallpapers.
//
// A wallpaper is produced in three steps: a random grayscale mask is created,
// a path is traced through every pixel of the mask, and each pixel's position
// along the path is mapped to a colour via a procedurally generated colourmap.
//
// The public API is minimal: call Generate with a width, height, and output
// name, optionally configuring the number of images or a random seed.
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

// Option configures a call to Generate.
type Option func(*config)

type config struct {
	count int
	seed  int64
}

// WithCount sets the number of images to generate (default 1).
func WithCount(n int) Option {
	return func(c *config) {
		c.count = n
	}
}

// WithSeed sets the random seed for reproducible output (default:
// time.Now().UnixNano()). The same seed and count always produce the same images.
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
		t := 0.0
		if pathLen > 1 {
			t = float64(i) / float64(pathLen-1)
		}
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
		defer func() {
			if cerr := f.Close(); cerr != nil {
				fmt.Fprintf(os.Stderr, "randwallpaper: warning: close failed: %v\n", cerr)
			}
		}()

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
