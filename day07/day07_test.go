package main

import (
	"testing"
)

func Test_part1(t *testing.T) {
	t.Run("part1", func(t *testing.T) {
		ic := parseInputToIntCode("input.txt")
		ans := 212460
		if got := part1(ic); got != ans {
			t.Errorf("part1() = %v, want %v", got, ans)
		}
	})
}

func Test_part2a(t *testing.T) {
	t.Run("part2a", func(t *testing.T) {
		ic := parseInputToIntCode("input.txt")
		ans := 21844737
		if got := part2a(ic); got != ans {
			t.Errorf("part2a() = %v, want %v", got, ans)
		}
	})
}
