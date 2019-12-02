package files

import (
	"bufio"
	"os"
	"strconv"
)

// Reads an entire file into memory and returns a slice of it's lines.
func ReadAllLines(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	s := bufio.NewScanner(f)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines, s.Err()
}

func ReadAllLinesAsInts(path string) ([]int, error) {
	longInts, err := ReadAllLinesAsInt64s(path)
	if err != nil {
		return nil, err
	}

	var ints []int
	for _, longInt := range longInts {
		ints = append(ints, int(longInt))
	}
	return ints, nil
}

func ReadAllLinesAsInt64s(path string) ([]int64, error) {
	lines, err := ReadAllLines(path)
	if err != nil {
		return nil, err
	}

	var ints []int64
	for _, line := range lines {
		n, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			return nil, err
		}
		ints = append(ints, n)
	}
	return ints, nil
}
