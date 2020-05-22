package main

import (
	"aoc2019/geom"
	ic "aoc2019/intcode"
	"log"
)

func main() {
	part1()
}

type robotState int

const (
	hitWall           = 0
	movedSuccessfully = 1
	foundOxygen       = 2
)

type moveCmd int

const (
	north = 1
	south = 2
	west  = 3
	east  = 4
)

type robot struct {
	visited map[geom.Point]robotState
	cmdIndexByPt map[geom.Point]int
	commands     []moveCmd
	ptsToVisit   chan geom.Point
	code         *ic.Mem
	curPt        geom.Point
}

func invertCmd(cmd moveCmd) moveCmd {
	switch cmd {
	case north:
		return south
	case south:
		return north
	case west:
		return east
	case east:
		return west
	default:
		panic("unreachable")
	}
}

func (r *robot) move(cmd moveCmd) robotState {
	var destPt geom.Point
	switch cmd {
	case north:
		destPt = r.curPt.North()
	case south:
		destPt = r.curPt.South()
	case east:
		destPt = r.curPt.East()
	case west:
		destPt = r.curPt.West()
	default:
		panic("unreachable")
	}

	s1 := <-r.code.State
	if s1 != ic.NeedInput {
		log.Fatalf("expected need input")
	}
	r.code.Input <- int(cmd)
	s2 := <-r.code.State
	if s2 != ic.HaveOutput {
		log.Fatalf("expected need output")
	}
	state := robotState(<-r.code.Output)

	switch state {
	case hitWall:
	case foundOxygen:
		fallthrough
	case movedSuccessfully:
		r.curPt = destPt
		r.commands = append(r.commands, cmd)
		r.cmdIndexByPt[destPt] = len(r.commands) - 1
	default:
		panic("unreachable")
	}
	return state
}


func (r *robot) moveToPt(pt geom.Point) robotState {
	var cmdIndex int
	var dir moveCmd
	if i, ok := r.cmdIndexByPt[pt]; ok {
		cmdIndex = i
	} else if i, ok := r.cmdIndexByPt[pt.North()]; ok {
		cmdIndex = i
		dir = south
	} else if i, ok := r.cmdIndexByPt[pt.South()]; ok {
		cmdIndex = i
		dir = north
	} else if i, ok := r.cmdIndexByPt[pt.East()]; ok {
		cmdIndex = i
		dir = west
	} else if i, ok := r.cmdIndexByPt[pt.West()]; ok {
		cmdIndex = i
		dir = east
	}

	for _, cmd := range r.commands[cmdIndex:] {
		r.move(invertCmd(cmd))
	}
	if dir != 0 {
		r.move(dir)
	}
}

func part1(origCode *ic.Mem) int {
	r := robot{
		visited:      make(map[geom.Point]robotState),
		cmdIndexByPt: make(map[geom.Point]int),
		commands:     make([]moveCmd, 0),
		ptsToVisit:   make(chan geom.Point, 300),
		code:         origCode.Clone(),
		curPt:        geom.NewPoint(0, 0),
	}
	r.ptsToVisit <- geom.NewPoint(0, 0)
	r.cmdIndexByPt[geom.NewPoint(0, 0)] = 0

	for {
		nextPt := <-r.ptsToVisit
		if _, ok := r.cmdIndexByPt[nextPt]; ok {
			continue
		}
		state := r.moveToPt(nextPt)
		switch state {
		case hitWall:
			continue
		case movedSuccessfully:
			r.ptsToVisit <- r.curPt.North()
			r.ptsToVisit <- r.curPt.South()
			r.ptsToVisit <- r.curPt.West()
			r.ptsToVisit <- r.curPt.East()


		case foundOxygen:
			return findShortestPath(r.visited, r.curPt)
		}
	}
}
func findShortestPath(states map[geom.Point]robotState, destPt geom.Point) int {
	return 2
}
