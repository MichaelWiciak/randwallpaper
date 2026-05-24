package randwallpaper

import (
	"image"
	"math"
	"math/rand/v2"
	"sort"
)

type path interface {
	Path() []image.Point
}

type basePath struct {
	mask        mask
	width       int
	height      int
	maxX, maxY int
}

func newBasePath(m mask) basePath {
	return basePath{
		mask:   m,
		width:  m.Width(),
		height: m.Height(),
		maxX:   m.Width() - 1,
		maxY:   m.Height() - 1,
	}
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func (bp *basePath) getSquareNeighborhood(cur image.Point, used map[image.Point]bool) []image.Point {
	x, y := cur.X, cur.Y
	cx := bp.width - 1 - x
	cy := bp.height - 1 - y
	maxRadius := max(max(x, cx), max(y, cy))

	for radius := 1; radius <= maxRadius; radius++ {
		seen := make(map[image.Point]struct{})

		for i := 0; i <= 2*radius; i++ {
			ny := clamp(y-radius+i, 0, bp.maxY)
			seen[image.Point{clamp(x-radius, 0, bp.maxX), ny}] = struct{}{}
			seen[image.Point{clamp(x+radius, 0, bp.maxX), ny}] = struct{}{}
		}
		for i := 1; i < 2*radius; i++ {
			nx := clamp(x-radius+i, 0, bp.maxX)
			seen[image.Point{nx, clamp(y-radius, 0, bp.maxY)}] = struct{}{}
			seen[image.Point{nx, clamp(y+radius, 0, bp.maxY)}] = struct{}{}
		}

		unvisited := make([]image.Point, 0, len(seen))
		for p := range seen {
			if !used[p] {
				unvisited = append(unvisited, p)
			}
		}
		if len(unvisited) > 0 {
			sort.Slice(unvisited, func(i, j int) bool {
				if unvisited[i].X != unvisited[j].X {
					return unvisited[i].X < unvisited[j].X
				}
				return unvisited[i].Y < unvisited[j].Y
			})
			return unvisited
		}
	}
	return nil
}

type epwtPath struct {
	basePath
	rng *rand.Rand
}

func newEPWTPath(m mask, rng *rand.Rand) *epwtPath {
	bp := newBasePath(m)
	return &epwtPath{basePath: bp, rng: rng}
}

func (p *epwtPath) Path() []image.Point {
	x, y := p.rng.IntN(p.width), p.rng.IntN(p.height)
	cur := image.Point{x, y}
	out := make([]image.Point, 0, p.width*p.height)
	out = append(out, cur)
	used := map[image.Point]bool{cur: true}

	for {
		neighbors := p.getSquareNeighborhood(cur, used)
		if len(neighbors) == 0 {
			break
		}
		curVal := p.mask.At(cur.X, cur.Y)
		best := neighbors[0]
		bestDiff := math.Abs(curVal - p.mask.At(best.X, best.Y))
		for _, n := range neighbors[1:] {
			diff := math.Abs(curVal - p.mask.At(n.X, n.Y))
			if diff < bestDiff {
				bestDiff = diff
				best = n
			}
		}
		out = append(out, best)
		used[best] = true
		cur = best
	}
	return out
}

type probabilisticPath struct {
	basePath
	rng *rand.Rand
}

func newProbabilisticPath(m mask, rng *rand.Rand) *probabilisticPath {
	bp := newBasePath(m)
	return &probabilisticPath{basePath: bp, rng: rng}
}

func (p *probabilisticPath) Path() []image.Point {
	x, y := p.rng.IntN(p.width), p.rng.IntN(p.height)
	cur := image.Point{x, y}
	out := make([]image.Point, 0, p.width*p.height)
	out = append(out, cur)
	used := map[image.Point]bool{cur: true}

	for {
		neighbors := p.getSquareNeighborhood(cur, used)
		if len(neighbors) == 0 {
			break
		}
		weights := make([]float64, len(neighbors))
		sum := 0.0
		for i, n := range neighbors {
			w := p.mask.At(n.X, n.Y)
			weights[i] = w
			sum += w
		}
		if sum == 0 {
			next := neighbors[p.rng.IntN(len(neighbors))]
			out = append(out, next)
			used[next] = true
			cur = next
			continue
		}
		target := p.rng.Float64() * sum
		cum := 0.0
		next := neighbors[len(neighbors)-1]
		for i, w := range weights {
			cum += w
			if target <= cum {
				next = neighbors[i]
				break
			}
		}
		out = append(out, next)
		used[next] = true
		cur = next
	}
	return out
}

func newPath(m mask, rng *rand.Rand) path {
	if rng.IntN(2) == 0 {
		return newEPWTPath(m, rng)
	}
	return newProbabilisticPath(m, rng)
}
