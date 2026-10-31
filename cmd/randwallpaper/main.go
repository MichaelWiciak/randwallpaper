package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime/debug"

	randwallpaper "github.com/MichaelWiciak/randwallpaper"
)

var version = "dev"

type parseError struct {
	error
}

func (e *parseError) Unwrap() error { return e.error }

func main() {
	err := run(os.Args[1:], os.Stdout, os.Stderr)
	if err == nil {
		return
	}

	var pe *parseError
	if errors.As(err, &pe) {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		os.Exit(2)
	}

	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

func run(args []string, stdout, stderr io.Writer) error {
	var width, height, count, seed int
	var output string
	var showVersion bool

	fs := flag.NewFlagSet("randwallpaper", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.IntVar(&width, "width", 1920, "image width")
	fs.IntVar(&width, "w", 1920, "image width (shorthand)")
	fs.IntVar(&height, "height", 1080, "image height")
	fs.IntVar(&height, "h", 1080, "image height (shorthand)")
	fs.StringVar(&output, "out", "wallpaper", "output file name (without extension)")
	fs.IntVar(&count, "count", 1, "number of images to generate")
	fs.IntVar(&seed, "seed", 0, "random seed (0 = random)")
	fs.BoolVar(&showVersion, "version", false, "print version and exit")

	fs.Usage = func() {
		fmt.Fprintf(stderr, "Usage: randwallpaper [options]\n\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		return &parseError{err}
	}

	if showVersion {
		fmt.Fprintln(stdout, resolveVersion())
		return nil
	}

	if width <= 0 || height <= 0 {
		return fmt.Errorf("width and height must be positive, got %dx%d", width, height)
	}
	if count <= 0 {
		return fmt.Errorf("count must be positive, got %d", count)
	}

	var opts []randwallpaper.Option
	if count > 1 {
		opts = append(opts, randwallpaper.WithCount(count))
	}
	if seed != 0 {
		opts = append(opts, randwallpaper.WithSeed(int64(seed)))
	}

	return randwallpaper.Generate(width, height, output, opts...)
}

func resolveVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return version
}
