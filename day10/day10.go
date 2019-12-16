package main

import (
	"aoc2019/files"
	"aoc2019/maths"
	"fmt"
	"math"
	"sort"
	"strings"
)

// Calculated by the smallest angle on a 40x40 grid.
const polarAngleEpsilon = 0.0007

func main() {
	rawLines := files.ReadAllLines("day10/input.txt")
	best, start := part1(rawLines)
	fmt.Printf("Part 1: best=%d, point=%s\n", best, start)

	asteroid200 := part2(start, rawLines, 200)

	size := len(rawLines) * 2
	p1 := toPolar(point{0, 0}, point{size, size})
	p2 := toPolar(point{0, 0}, point{size, size - 1})
	fmt.Printf("min angle: %s - %s: %f\n", p1, p2, math.Abs(p1.angle-p2.angle))

	fmt.Printf("Part 2: asteroidPolar=%s\n", asteroid200)
}

type point struct{ x, y int }

type polar struct{ angle, radius float64 }

func (p point) String() string {
	return fmt.Sprintf("point{%d %d}", p.x, p.y)
}

func (p polar) String() string {
	return fmt.Sprintf("polar{%f (%.1f deg), %f}", p.angle, p.angle*180/math.Pi, p.radius)
}

func (p polar) toCartesian(start point) point {
	x := math.Sin(p.angle) * p.radius
	y := math.Cos(p.angle) * p.radius * -1
	return point{start.x + int(math.Round(x)), start.y + int(math.Round(y))}
}

func ptLeftOf(pt1, pt2 point) bool {
	if pt1.x < pt2.x {
		return true
	} else if pt1.x == pt2.x {
		return pt1.y < pt2.y
	} else {
		return false
	}
}

type line struct {
	dx, dy     int
	yIntercept float64
	// Is 0 unless a vertical line
	xIntercept int
}

func (l line) String() string {
	return fmt.Sprintf("line{dy=%d dx=%d y=%f x=%d}", l.dy, l.dx, l.yIntercept, l.xIntercept)
}

func newLine(pt1, pt2 point) line {
	dx := pt1.x - pt2.x
	dy := pt1.y - pt2.y
	yIntercept := float64(0)
	xIntercept := 0
	if dy < 0 {
		dx *= -1
		dy *= -1
	}
	if dy == 0 {
		// Normalize dx.
		dx = 1
		yIntercept = float64(pt1.y)
	} else if dx == 0 {
		// Normalize dy.
		dy = 1
		xIntercept = pt1.x
	} else {
		// Reduce to proper fraction.
		target := maths.MaxInt(maths.AbsInt(dx), maths.AbsInt(dy))
		n := 2
		for n <= target {
			for dx%n == 0 && dy%n == 0 {
				dx /= n
				dy /= n
			}
			n += 1
		}
		// Formula:
		// y = (dy/dx)x + b
		// b = y - (dy/dx)x
		yIntercept = float64(pt1.y) - float64(dy)/float64(dx)*float64(pt1.x)
	}

	return line{
		dx:         dx,
		dy:         dy,
		yIntercept: yIntercept,
		xIntercept: xIntercept,
	}
}

func part1(rawLines []string) (int, point) {

	points := make([]point, 0, len(rawLines)*len(rawLines[0]))
	for y, line := range rawLines {
		for x, pos := range strings.Split(line, "") {
			if pos == "#" {
				points = append(points, point{x, y})
			}
		}
	}

	ptsOnLine := make(map[line][]point)
	lineHasPts := make(map[line]map[point]bool)
	for i, pt1 := range points {
		for j := i + 1; j < len(points); j++ {
			pt2 := points[j]
			l := newLine(pt1, pt2)
			if lineHasPts[l] == nil {
				lineHasPts[l] = make(map[point]bool)
			}
			if _, ok := lineHasPts[l][pt1]; !ok {
				ptsOnLine[l] = append(ptsOnLine[l], pt1)
				lineHasPts[l][pt1] = true
			}
			if _, ok := lineHasPts[l][pt2]; !ok {
				ptsOnLine[l] = append(ptsOnLine[l], pt2)
				lineHasPts[l][pt2] = true
			}
		}
	}

	for _, pts := range ptsOnLine {
		lessThan := func(i, j int) bool { return ptLeftOf(pts[i], pts[j]) }
		sort.Slice(pts, lessThan)
	}

	best := 0
	bestPt := point{math.MinInt64, math.MinInt64}
	for _, p := range points {
		// Assume we can see all points
		seen := make(map[point]bool)
		for _, p2 := range points {
			seen[p2] = true
		}
		delete(seen, p)

		for l, pts := range ptsOnLine {
			if _, ok := lineHasPts[l][p]; !ok {
				// If the current point isn't in the line, we can see all points
				// in this line.
				continue
			}

			// The index of the current point in this line of points
			idx := -1
			for i, pt := range pts {
				if pt == p {
					idx = i
					break
				}
			}

			// Remove all non adjacent points since they're blocked.
			for i := 0; i < idx-1; i++ {
				delete(seen, pts[i])
			}
			for i := idx + 2; i < len(pts); i++ {
				delete(seen, pts[i])
			}

		}
		if len(seen) > best {
			best = len(seen)
			bestPt = p
		}
	}
	return best, bestPt
}

// Returns the polar coordinates in [0, 2*pi) with 0 = straight up.
//
func toPolar(start, dest point) polar {
	dx := float64(dest.x - start.x)
	dy := float64(start.y - dest.y)
	if dy == 0 {
		if dx == 0 {
			panic("had same point")
		} else if dx < 0 {
			return polar{3 * math.Pi / 2, math.Abs(dx)}
		} else {
			return polar{math.Pi / 2, dx}
		}
	} else {
		angle := math.Atan2(dx, dy)
		if angle < 0 {
			angle += 2 * math.Pi
		}
		radius := math.Sqrt(math.Pow(dx, 2) + math.Pow(dy, 2))
		return polar{angle, radius}
	}
}

func part2(start point, rawLines []string, numAsteroids int) point {
	allPoints := parseLines(rawLines)
	points := make([]point, 0, len(allPoints)-1)
	for _, pt := range allPoints {
		if pt != start {
			points = append(points, pt)
		}
	}

	polars := make([]polar, 0, len(points))
	for _, dest := range points {
		p := toPolar(start, dest)
		polars = append(polars, p)
	}
	sort.SliceStable(polars, func(i, j int) bool {

		if polars[i].angle < polars[j].angle {
			return true
		} else if polars[i].angle == polars[j].angle {
			return polars[i].radius < polars[j].radius
		} else {
			return false
		}
	})

	lastAngle := math.MaxFloat64
	count := 0
	destroyed := make([]bool, len(polars))
	for i := 0; count < numAsteroids; i = (i + 1) % len(polars) {
		// Reset lastAngle if we make it all the way around.
		if i == 0 {
			lastAngle = math.MaxFloat64
		}
		if destroyed[i] {
			continue
		}
		p := polars[i]
		if math.Abs(p.angle-lastAngle) < polarAngleEpsilon {
			if p.angle != lastAngle {
				fmt.Printf("  unequal but same %s %f\n", p, lastAngle)
			}
			fmt.Printf("Skipping asteroid   i=%d count=%d at %s, %s, lastAngle=%f\n", i, count, p, p.toCartesian(start), lastAngle)
			continue
		}
		count++
		fmt.Printf("Destroying asteroid i=%d count=%d at %s, %s lastAngle=%f\n", i, count, p, p.toCartesian(start), lastAngle)
		lastAngle = p.angle
		destroyed[i] = true
		if count == numAsteroids {
			return polars[i].toCartesian(start)
		}
	}
	panic("exited loop")
}

func parseLines(rawLines []string) []point {
	points := make([]point, 0, len(rawLines)*len(rawLines[0]))
	for y, line := range rawLines {
		for x, pos := range strings.Split(line, "") {
			if pos == "#" {
				points = append(points, point{x, y})
			}
		}
	}
	return points
}
