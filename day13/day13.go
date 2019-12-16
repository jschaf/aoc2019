package main

import (
	"aoc2019/files"
	"aoc2019/geom"
	"aoc2019/intcode"
	"fmt"
)

func main() {
	line := files.ReadFirstLine("day13/input.txt")
	mem := intcode.ParseLine(line)
	s := part1(mem)
	blocks := 0
	for _, runes := range s.buf {
		for _, c := range runes {
			if c == blockChar {
				blocks++
			}
		}
	}
	fmt.Printf("Part1: num blocks: %d\n", blocks)
	s.drawScreen()

	s2 := part2(mem)
	fmt.Printf("Part2: high score: %d\n", s2.highScore)
}

type tileId = int

const (
	emptyTile  tileId = 0
	wallTile   tileId = 1
	blockTile  tileId = 2
	paddleTile tileId = 3
	ballTile   tileId = 4
)

const (
	wallChar   = '|'
	blockChar  = 'B'
	paddleChar = '_'
	ballChar   = '0'
)

func part1(mem *intcode.Mem) *screen {
	code := mem.Clone()
	go code.Run()

	screen := newScreen()

	for {
		select {
		case s := <-code.State:
			switch s {
			case intcode.HaveOutput:
				x := <-code.Output
				<-code.State
				y := <-code.Output
				<-code.State
				tileId := <-code.Output
				screen.drawTile(x, y, tileId)
			case intcode.NeedInput:
				panic("intcode needs input")
			case intcode.Halted:
				return screen
			default:
				panic("unexpected state")
			}
		}
	}
}

type screen struct {
	buf                [][]rune
	highScore          int
	paddlePos, ballPos geom.Point
}

func newScreen() *screen {
	size := 40 // determined by inspection
	b := make([][]rune, size)
	for i := range b {
		b[i] = make([]rune, size)
		for j := range b[i] {
			b[i][j] = ' '
		}
	}
	return &screen{buf: b}
}

func (s *screen) drawTile(x, y int, t tileId) {
	switch t {
	case emptyTile:
		s.buf[y][x] = ' '
	case wallTile:
		s.buf[y][x] = wallChar
	case blockTile:
		s.buf[y][x] = blockChar
	case paddleTile:
		s.paddlePos = geom.NewPoint(x, y)
		s.buf[y][x] = paddleChar
	case ballTile:
		s.ballPos = geom.NewPoint(x, y)
		s.buf[y][x] = ballChar
	default:
		panic("unexpected tileId")
	}
}

func (s *screen) drawScreen() {
	for _, runes := range s.buf {
		fmt.Println(string(runes))
	}
}

const (
	joystickNeutral = 0
	joystickLeft    = -1
	joystickRight   = 1
)

func part2(mem *intcode.Mem) *screen {
	code := mem.Clone()
	numQuarters := 2
	code.Set(0, numQuarters)
	go code.Run()

	screen := newScreen()

	for {
		select {
		case s := <-code.State:
			switch s {
			case intcode.HaveOutput:
				x := <-code.Output
				<-code.State
				y := <-code.Output
				<-code.State
				tileId := <-code.Output
				if x == -1 && y == 0 {
					screen.highScore = tileId
					fmt.Printf("updated high score: %d\n", tileId)
				} else {
					screen.drawTile(x, y, tileId)
				}
			case intcode.NeedInput:
				p := screen.paddlePos
				b := screen.ballPos
				if b.X < p.X {
					code.Input <- joystickLeft
				} else if b.X == p.X {
					code.Input <- joystickNeutral
				} else {
					code.Input <- joystickRight
				}
			case intcode.Halted:
				return screen
			default:
				panic("unexpected state")
			}
		}
	}
}
