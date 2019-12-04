package main

import (
	"aoc2019/files"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func main() {
	lines := files.ReadAllLines("day04/input.txt")
	split := strings.Split(lines[0], "-")
	lo, _ := strconv.Atoi(split[0])
	hi, _ := strconv.Atoi(split[1])

	fmt.Println(countMatches(lo, hi))
	fmt.Printf("part2: %d\n", countMatches2(lo, hi))
}

func countMatches(lo int, hi int) int {
	numMatching := 0
	for i := lo; i <= hi; i += 1 {
		if isMatch(i) {
			numMatching += 1
		}
	}
	return numMatching
}

func isMatch(i int) bool {
	lastDigit := math.MaxInt64
	hadAdjacent := false
	incrDigits := true
	n := i
	for n > 0 {
		d := n % 10
		hadAdjacent = hadAdjacent || d == lastDigit
		incrDigits = incrDigits && d <= lastDigit
		lastDigit = d
		n /= 10
	}
	return hadAdjacent && incrDigits
}

func countMatches2(lo int, hi int) int {
	numMatching := 0
	for x := lo; x <= hi; x += 1 {
		if isMatch2(x) {
			numMatching += 1
		}
	}
	return numMatching
}

func isMatch2(x int) bool {
	prevD := math.MaxInt64
	prevPrevD := math.MaxInt64 - 1
	hasRun2 := false
	isLongRun := false
	isDigitsDecreasing := true
	n := x

	for n > 0 {
		d := n % 10

		// We have a run of 2 digits if the previous 2 digits match.
		curRun2 := (d != prevD) && (prevD == prevPrevD) && !isLongRun
		hasRun2 = hasRun2 || curRun2
		isDigitsDecreasing = isDigitsDecreasing && d <= prevD

		// Update
		isLongRun = d == prevD && d == prevPrevD
		prevPrevD = prevD
		prevD = d
		n /= 10
	}

	// We haven't checked the 2 most significant digits yet.
	lastRun2 := (prevD == prevPrevD) && !isLongRun
	hasRun2 = hasRun2 || lastRun2
	return hasRun2 && isDigitsDecreasing
}
