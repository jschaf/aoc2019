package files

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// Reads an entire file into memory and returns a slice of it's lines.
func ReadAllLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines
}

func ReadFirstLine(path string) string {
	lines := ReadAllLines(path)
	if len(lines) == 0 {
		log.Fatalf("expected 1 line but had 0 from file: %s", path)
	}
	if len(lines) > 1 {
		log.Fatalf("expected 1 line but had %d lines from file: %s", len(lines), path)
	}
	return lines[0]
}

func ReadAllLinesAsInts(path string) []int {
	longInts := ReadAllLinesAsInt64s(path)

	var ints []int
	for _, longInt := range longInts {
		ints = append(ints, int(longInt))
	}
	return ints
}

func ReadAllLinesAsInt64s(path string) []int64 {
	lines := ReadAllLines(path)

	var ints []int64
	for _, line := range lines {
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		ints = append(ints, n)
	}
	return ints
}
