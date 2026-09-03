package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	randwallpaper "github.com/MichaelWiciak/randwallpaper"
)

func main() {
	var width, height, count, seed int
	var output string
	var showVersion bool

	flag.IntVar(&width, "width", 1920, "image width")
	flag.IntVar(&width, "w", 1920, "image width (shorthand)")
	flag.IntVar(&height, "height", 1080, "image height")
	flag.IntVar(&height, "h", 1080, "image height (shorthand)")
	flag.StringVar(&output, "out", "wallpaper", "output file name (without extension)")
	flag.IntVar(&count, "count", 1, "number of images to generate")
	flag.IntVar(&seed, "seed", 0, "random seed (0 = random)")
	flag.BoolVar(&showVersion, "version", false, "print version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		version := "dev"
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" {
			version = info.Main.Version
		}
		fmt.Println(version)
		return
	}

	if width <= 0 || height <= 0 {
		log.Fatalf("width and height must be positive, got %dx%d", width, height)
	}
	if count <= 0 {
		log.Fatalf("count must be positive, got %d", count)
	}

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
