package randwallpaper

import (
	"math"
	"runtime"
	"sync"
	"sync/atomic"
)

func gaussianBlur(src []uint8, width, height int, sigma float64) []uint8 {
	dst := make([]uint8, len(src))
	if width <= 0 || height <= 0 || sigma <= 0 {
		copy(dst, src)
		return dst
	}

	radius := int(math.Ceil(sigma * 3.0))
	kernel := make([]float64, radius+1)
	for i := range kernel {
		x := float64(i)
		kernel[i] = math.Exp(-(x*x)/(2*sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))
	}

	tmp := make([]uint8, len(src))
	parallelFor(height, func(y int) {
		for x := range width {
			lo := max(0, x-radius)
			hi := min(width-1, x+radius)
			var r, a float64
			for ix := lo; ix <= hi; ix++ {
				wa := 255 * kernel[absInt(x-ix)]
				r += float64(src[y*width+ix]) * wa
				a += wa
			}
			if a != 0 {
				tmp[y*width+x] = clampByte(r * (1 / a))
			}
		}
	})

	parallelFor(height, func(y int) {
		lo := max(0, y-radius)
		hi := min(height-1, y+radius)
		for x := range width {
			var r, a float64
			for iy := lo; iy <= hi; iy++ {
				wa := 255 * kernel[absInt(y-iy)]
				r += float64(tmp[iy*width+x]) * wa
				a += wa
			}
			if a != 0 {
				dst[y*width+x] = clampByte(r * (1 / a))
			}
		}
	})

	return dst
}

func parallelFor(n int, fn func(int)) {
	workers := min(runtime.GOMAXPROCS(0), n)
	if workers <= 1 {
		for i := range n {
			fn(i)
		}
		return
	}
	var next atomic.Int64
	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			for {
				i := int(next.Add(1)) - 1
				if i >= n {
					return
				}
				fn(i)
			}
		}()
	}
	wg.Wait()
}

func absInt(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func clampByte(v float64) uint8 {
	n := int64(v + 0.5)
	if n > 255 {
		return 255
	}
	if n > 0 {
		return uint8(n)
	}
	return 0
}
