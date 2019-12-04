package main

import (
	"aoc2019/files"
	"aoc2019/lines"
	"aoc2019/maths"
	"fmt"
	"log"
	"math"
	"strconv"
)

func main() {
	ls := files.ReadAllLines("day03/input.txt")
	if len(ls) != 2 {
		panic("expected 2 lines")
	}
	var grids []grid

	for _, l := range ls {
		directions := lines.ParseCommaSeparatedStrings(l)
		grids = append(grids, makeGrid(directions))
	}

	g1 := grids[0]
	g2 := grids[1]
	best := math.MaxInt32
	bestTiming := math.MaxInt32
	for x, ys := range g1 {
		for y, d1 := range ys {
			d2 := g2[x][y]
			if d2 > 0 {
				best = maths.MinInt(maths.AbsInt(x)+maths.AbsInt(y), best)
				bestTiming = maths.MinInt(d1+d2, bestTiming)
			}
		}
	}
	fmt.Printf("Part 1: %d\n", best)
	fmt.Printf("Part 2: %d\n", bestTiming)
}

type gridCol map[int]int
type grid map[int]gridCol

func set(g grid, x, y, new int) {
	if g[x] == nil {
		g[x] = make(map[int]int)
	}
	old := g[x][y]
	if old == 0 {
		g[x][y] = new
	} else {
		g[x][y] = maths.MinInt(new, old)
	}
}

func makeGrid(directions []string) grid {
	g := make(grid)
	x, y, steps := 0, 0, 0
	for _, dir := range directions {
		distRaw, _ := strconv.Atoi(dir[1:])
		dist := distRaw
		dx, dy := 0, 0
		switch rune(dir[0]) {
		case 'U':
			dy = 1
		case 'R':
			dx = 1
		case 'D':
			dy = -1
		case 'L':
			dx = -1
		default:
			log.Fatalf("unexpected direction: %s", string(dir[0]))
		}
		for n := 0; n < dist; n++ {
			x += dx
			y += dy
			steps += 1
			set(g, x, y, steps)
		}
	}
	return g
}
