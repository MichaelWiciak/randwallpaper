package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	randwallpaper "github.com/MichaelWiciak/randwallpaper"
)

func main() {
	var width, height, count, seed int
	var output string

	flag.IntVar(&width, "width", 1920, "image width")
	flag.IntVar(&width, "w", 1920, "image width (shorthand)")
	flag.IntVar(&height, "height", 1080, "image height")
	flag.IntVar(&height, "h", 1080, "image height (shorthand)")
	flag.StringVar(&output, "out", "wallpaper", "output file name (without extension)")
	flag.IntVar(&count, "count", 1, "number of images to generate")
	flag.IntVar(&seed, "seed", 0, "random seed (0 = random)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	var opts []randwallpaper.Option
	if count > 1 {
		opts = append(opts, randwallpaper.WithCount(count))
	}
	if seed != 0 {
		opts = append(opts, randwallpaper.WithSeed(int64(seed)))
	}

	if err := randwallpaper.Generate(width, height, output, opts...); err != nil {
		log.Fatal(err)
	}
}
