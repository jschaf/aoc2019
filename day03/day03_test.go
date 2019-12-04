package main

import (
	"reflect"
	"testing"
)

type point struct{ x, y int }

func newGrid(pts ...point) grid {
	grid := make(grid)
	for i, pt := range pts {
		if grid[pt.x] == nil {
			grid[pt.x] = make(gridCol)
		}
		grid[pt.x][pt.y] = i + 1
	}
	return grid
}

func dirs(ds ...string) []string {
	var dirs []string
	for _, d := range ds {
		dirs = append(dirs, d)
	}
	return dirs
}

func Test_makeGrid(t *testing.T) {
	tests := []struct {
		name       string
		directions []string
		want       grid
	}{
		{"U1", dirs("U1"), newGrid(point{0, 1})},
		{"U1,R2", dirs("U1", "R2"),
			newGrid(point{0, 1}, point{1, 1}, point{2, 1}),
		},
		{"U1,L2", dirs("U1", "L2"),
			newGrid(point{0, 1}, point{-1, 1}, point{-2, 1}),
		},
		{"D1,L2", dirs("D1", "L2"),
			newGrid(point{0, -1}, point{-1, -1}, point{-2, -1}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeGrid(tt.directions); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("makeGrid() = %v, want %v", got, tt.want)
			}
		})
	}
}
