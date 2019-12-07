package main

import (
	"aoc2019/files"
	"fmt"
	"log"
	"math"
	"strings"
)

type orbit struct {
	id     string
	parent *orbit
}

func main() {
	lines := files.ReadAllLines("day06/input.txt")
	allOrbits := make(map[string]*orbit)
	for _, line := range lines {
		os := strings.Split(line, ")")
		pId := os[0]
		cId := os[1]
		if _, ok := allOrbits[pId]; !ok {
			allOrbits[pId] = &orbit{pId, nil}
		}

		p := allOrbits[pId]
		if c, ok := allOrbits[cId]; ok {
			c.parent = p
		} else {
			allOrbits[cId] = &orbit{cId, p}
		}
	}

	total := 0
	for _, o := range allOrbits {
		if o.parent == nil && o.id != "COM" {
			log.Fatalf("Only COM may not have a parent but had %s", o.id)
		}
		for o.parent != nil {
			total += 1
			o = o.parent
		}
	}

	log.Printf("Part 1 total orbits: %d", total)

	// Get SAN parents as a map[id]int start at 0
	// Get YOU parents as a map[id]int
	// Pick smallest intersection
	// Add them, subtract 1?
	sanMap := make(map[string]int)
	san := allOrbits["SAN"]
	sp := san.parent
	sl := 0
	for sp != nil {
		sanMap[sp.id] = sl
		sl += 1
		sp = sp.parent
	}

	youMap := make(map[string]int)
	you := allOrbits["YOU"]
	yp := you.parent
	yl := 0
	for yp != nil {
		youMap[yp.id] = yl
		yl += 1
		yp = yp.parent
	}

	shortest := math.MaxInt64
	for id, length := range sanMap {
		if l, ok := youMap[id]; ok {
			if length+l < shortest {
				shortest = length + l
			}
		}
	}
	fmt.Printf("Part 2 num transfers: %d", shortest)

}
