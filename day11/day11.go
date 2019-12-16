package main

import (
	"aoc2019/files"
	"aoc2019/intcode"
	"aoc2019/maths"
	"fmt"
	"log"
)

func main() {
	line := files.ReadFirstLine("day11/input.txt")
	code := intcode.ParseLine(line)
	numTiles := len(runPaintRobot(code, black).colors)
	fmt.Printf("Part1: %d tiles painted\n", numTiles)

	hull := part2(code)
	fmt.Println("Part2 hull:")
	for _, runes := range hull {
		fmt.Println(string(runes))
	}
}

type point struct{ x, y int }

type orientation int

const (
	up = iota
	left
	right
	down
)

const (
	black = 0
	white = 1
)

const (
	turnLeft  = 0
	turnRight = 1
)

type hullColor int

type paintState struct {
	// The direction we currently face.
	dir orientation
	// The current position.
	pos    point
	colors map[point]hullColor
}

func (ps *paintState) turnLeft() {
	switch ps.dir {
	case up:
		ps.dir = left
	case left:
		ps.dir = down
	case down:
		ps.dir = right
	case right:
		ps.dir = up
	default:
		panic("unexpected direction")
	}
}

func (ps *paintState) turnRight() {
	switch ps.dir {
	case up:
		ps.dir = right
	case left:
		ps.dir = up
	case down:
		ps.dir = left
	case right:
		ps.dir = down
	default:
		panic("unexpected direction")
	}
}

func (ps *paintState) advance() {
	switch ps.dir {
	case up:
		ps.pos = point{ps.pos.x, ps.pos.y - 1}
	case left:
		ps.pos = point{ps.pos.x - 1, ps.pos.y}
	case down:
		ps.pos = point{ps.pos.x, ps.pos.y + 1}
	case right:
		ps.pos = point{ps.pos.x + 1, ps.pos.y}
	default:
		panic(fmt.Sprintf("invalid direction: %d", ps.dir))
	}
}

func (ps *paintState) paint(c hullColor) {
	ps.colors[ps.pos] = c
}

func runPaintRobot(mem *intcode.Mem, startColor int) paintState {
	code := mem.Clone()
	code.SeedInput([]int{startColor})
	go code.Run()

	state := paintState{
		dir:    up,
		pos:    point{0, 0},
		colors: map[point]hullColor{},
	}
	for {
		select {
		case s := <-code.State:
			switch s {
			case intcode.Halted:
				return state

			case intcode.NeedInput:
				code.Input <- int(state.colors[state.pos])

			case intcode.HaveOutput:
				c := hullColor(<-code.Output)
				s2 := <-code.State
				if s2 != intcode.HaveOutput {
					log.Fatalf("expected 2 outputs in a row")
				}
				turn := <-code.Output
				state.paint(c)
				switch turn {
				case turnLeft:
					state.turnLeft()
				case turnRight:
					state.turnRight()
				default:
					panic(fmt.Sprintf("unexpected turn direction %d", turn))
				}
				state.advance()
			}
		}
	}
}

func part2(mem *intcode.Mem) [][]rune {
	state := runPaintRobot(mem, white)
	xLo, xHi, yLo, yHi := 0, 0, 0, 0
	for pt := range state.colors {
		xLo = maths.MinInt(xLo, pt.x)
		xHi = maths.MaxInt(xHi, pt.x)
		yLo = maths.MinInt(yLo, pt.y)
		yHi = maths.MaxInt(yHi, pt.y)
	}
	xOffset := maths.MaxInt(0, xLo*-1)
	yOffset := maths.MaxInt(0, yLo*-1)

	hull := make([][]rune, yHi+yOffset+1)
	for i := range hull {
		hull[i] = make([]rune, xHi+xOffset+1)
		for j := range hull[i] {
			hull[i][j] = '.'
		}
	}

	for pt, color := range state.colors {
		x := pt.x + xOffset
		y := pt.y + yOffset
		switch color {
		case black:
			hull[y][x] = ' '
		case white:
			hull[y][x] = '#'
		default:
			panic("unexpected color")
		}
	}

	return hull
}
