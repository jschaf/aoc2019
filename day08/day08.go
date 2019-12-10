package main

import (
	"aoc2019/files"
	"aoc2019/lines"
	"fmt"
	"log"
	"math"
)

func main() {
	line := files.ReadFirstLine("day08/input.txt")
	ints := lines.ParseAdjacentSingleDigitInts(line)
	part1(ints)

	img := part2(ints, width, height)
	for _, colors := range img {
		for _, c := range colors {
			switch c {
			case Black:
				fmt.Printf(" ")
			case White:
				fmt.Printf("#")
			case Transparent:
				fmt.Printf(" ")
			default:
				log.Fatalf("unknown color: %d", c)
			}
		}
		fmt.Println()
	}
}

const (
	width  = 25
	height = 6
)

func part1(ints []int) {
	digitCountsByLayer := make(map[int]map[int]int)
	numLayers := len(ints) / (width * height)
	for i := 0; i < numLayers; i++ {
		digitCountsByLayer[i] = make(map[int]int)
	}
	for i, d := range ints {
		layer := i / (width * height)
		digitCountsByLayer[layer][d] += 1
	}

	least0Layer, least0Count := -1, math.MaxInt64
	for layer, countsByDigit := range digitCountsByLayer {
		numZeroes := countsByDigit[0]
		if numZeroes < least0Count {
			least0Layer = layer
			least0Count = numZeroes
		}
	}

	numOnes := digitCountsByLayer[least0Layer][1]
	numTwos := digitCountsByLayer[least0Layer][2]
	fmt.Printf("Part 1: num 1 * num 2: %d\n", numOnes*numTwos)
}

const (
	Black       = 0
	White       = 1
	Transparent = 2
)

func part2(ints []int, width, height int) [][]int {
	numLayers := len(ints) / (width * height)
	img := make([][]int, height)

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			c := Transparent
			for j := 0; j < numLayers; j++ {
				layerStart := (width * height) * j
				offset := (row * width) + col
				c2 := ints[layerStart+offset]
				if c2 == Black {
					c = Black
					break
				} else if c2 == White {
					c = White
					break
				}
			}
			img[row] = append(img[row], c)
		}
	}
	return img
}
