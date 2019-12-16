package main

import (
	"aoc2019/files"
	"aoc2019/intcode"
	"fmt"
	"log"
)

func main() {
	line := files.ReadFirstLine("day11/input.txt")
	code := intcode.ParseLine(line)
	numTiles := part1(code)
	fmt.Printf("Part1: %d tiles painted\n", numTiles)
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
	}
}

func (ps *paintState) advance() {
	switch ps.dir {
	case up:
		ps.pos = point{ps.pos.x, ps.pos.y + 1}
	case left:
		ps.pos = point{ps.pos.x - 1, ps.pos.y}
	case down:
		ps.pos = point{ps.pos.x, ps.pos.y - 1}
	case right:
		ps.pos = point{ps.pos.x + 1, ps.pos.y}
	default:
		panic(fmt.Sprintf("invalid direction: %d", ps.dir))
	}
}

func (ps *paintState) paint(c hullColor) {
	ps.colors[ps.pos] = c
}

func part1(code *intcode.Mem) int {
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
				return len(state.colors)

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
