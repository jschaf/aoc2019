package main

import (
	"testing"
)

func Test_part1(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
		ic := parseInputToIntCode("input.txt")
		ans := 2316632620
		if got := part1(ic); got != ans {
			t.Errorf("part1() = %v, want %v", got, ans)
		}
	})
}

func Test_part2(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
		ic := parseInputToIntCode("input.txt")
		ans := 78869
		if got := part2(ic); got != ans {
			t.Errorf("part2() = %v, want %v", got, ans)
		}
	})
}
