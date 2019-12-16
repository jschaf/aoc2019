package main

import (
	"aoc2019/files"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func Test_newLine(t *testing.T) {
	tests := []struct {
		name     string
		pt1, pt2 point
		want     line
	}{
		{"simple 0,0 -> 1,1", point{0, 0}, point{1, 1}, line{1, 1, 0, 0}},
		{"simple 1,1 -> 0,0", point{1, 1}, point{0, 0}, line{1, 1, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newLine(tt.pt1, tt.pt2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newLine_allPointsInLineHaveSameLine(t *testing.T) {
	tests := []struct {
		name   string
		points []point
		want   line
	}{
		{
			"1/1 at y=2",
			[]point{{0, 2}, {1, 3}, {-1, 1}, {-5, -3}, {19, 21}},
			line{1, 1, 2, 0},
		},
		{
			"5/3 at y=2",
			[]point{{0, 2}, {5, 5}, {10, 8}, {15, 11}, {20, 14}},
			line{5, 3, 2, 0},
		},
		{
			"0/1 at x=2",
			[]point{{2, -1}, {2, 0}, {2, 1}, {2, 11}, {2, 14}},
			line{0, 1, 0, 2},
		},
		{
			"1,0 at y=3",
			[]point{{-1, 3}, {0, 3}, {1, 3}, {2, 3}, {51, 3}},
			line{1, 0, 3, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, pt1 := range tt.points {
				for j := i + 1; j < len(tt.points); j++ {
					pt2 := tt.points[j]
					if got := newLine(pt1, pt2); !reflect.DeepEqual(got, tt.want) {
						t.Errorf("newLine(%s, %s) = %s, want %s", pt1, pt2, got, tt.want)
					}
					if got := newLine(pt2, pt1); !reflect.DeepEqual(got, tt.want) {
						t.Errorf("newLine(%s, %s) = %s, want %s", pt2, pt1, got, tt.want)
					}
				}
			}
		})

	}
}

func Test_part1(t *testing.T) {
	tests := []struct {
		name     string
		rawLines []string
		want     int
	}{
		{"1x1, all asteroids", []string{"#"}, 0},
		{"1x1, no asteroids", []string{" "}, 0},
		{"1x2, all asteroids", []string{"##"}, 1},
		{"1x3, all asteroids", []string{"###"}, 2},
		{"1x3, some asteroids", []string{"#.#"}, 1},
		{"1x3, some asteroids", []string{"#.#"}, 1},
		{"2x3, some asteroids", []string{
			"#.#",
			".#.",
		}, 2},
		{"3x3, all asteroids", []string{
			"###",
			"###",
			"###",
		}, 8},
		{"3x3, some asteroids", []string{
			"###",
			"#.#",
			"###",
		}, 7},
		{"5x5, some asteroids", []string{
			".#..#",
			".....",
			"#####",
			"....#",
			"...##",
		}, 8},
		{"8x8, some asteroids", []string{
			"......#.#.",
			"#..#.#....",
			"..#######.",
			".#.#.###..",
			".#..#.....",
			"..#....#.#",
			"#..#....#.",
			".##.#..###",
			"##...#..#.",
			".#....####",
		}, 33},
		{"10x10, example 1", []string{
			"#.#...#.#.",
			".###....#.",
			".#....#...",
			"##.#.#.#.#",
			"....#.#.#.",
			".##..###.#",
			"..#...##..",
			"..##....##",
			"......#...",
			".####.###.",
		}, 35},
		{"huge, example 1", []string{
			".#..##.###...#######",
			"##.############..##.",
			".#.######.########.#",
			".###.#######.####.#.",
			"#####.##.#.##.###.##",
			"..#####..#.#########",
			"####################",
			"#.####....###.#.#.##",
			"##.#################",
			"#####.##.###..####..",
			"..######..##.#######",
			"####.##.####...##..#",
			".#####..#.######.###",
			"##...#.##########...",
			"#.##########.#######",
			".####.#.###.###.#.##",
			"....##.##.###..#####",
			".#.#.###########.###",
			"#.#.#.#####.####.###",
			"###.##.####.##.#..##",
		}, 210},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := part1(tt.rawLines); got != tt.want {
				t.Errorf("part1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPart1(b *testing.B) {
	rawLines := files.ReadAllLines("input.txt")
	for n := 0; n < b.N; n++ {
		p1, _ := part1(rawLines)
		fmt.Printf("benchmark: %d", p1)
	}

}

func Test_toPolar(t *testing.T) {
	sqrt2 := math.Sqrt(2)
	tests := []struct {
		name        string
		start, dest point
		want        polar
	}{
		{"dx= 0, dy= 1", point{3, 2}, point{3, 3}, polar{math.Pi, 1}},
		{"dx= 0, dy= 2", point{3, 2}, point{3, 4}, polar{math.Pi, 2}},
		{"dx= 1, dy= 1", point{3, 2}, point{4, 3}, polar{3 * math.Pi / 4, sqrt2}},
		{"dx= 1, dy= 0", point{3, 2}, point{4, 2}, polar{2 * math.Pi / 4, 1}},
		{"dx= 1, dy=-1", point{3, 2}, point{4, 1}, polar{1 * math.Pi / 4, sqrt2}},
		{"dx= 0, dy=-1", point{3, 2}, point{3, 1}, polar{0, 1}},
		{"dx=-1, dy=-1", point{3, 2}, point{2, 1}, polar{7 * math.Pi / 4, sqrt2}},
		{"dx=-1, dy= 0", point{3, 2}, point{2, 2}, polar{6 * math.Pi / 4, 1}},
		{"dx=-1, dy= 1", point{3, 2}, point{2, 3}, polar{5 * math.Pi / 4, sqrt2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toPolar(tt.start, tt.dest); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toPolar() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_polar_toCartesian(t *testing.T) {
	sqrt2 := math.Sqrt(2)
	tests := []struct {
		name  string
		p     polar
		start point
		want  point
	}{
		{"polar(180,1) from (0,0)", polar{math.Pi, 1}, point{0, 0}, point{0, 1}},
		{"polar(0,1) from (1,1)", polar{0, 1}, point{1, 1}, point{1, 0}},
		{"polar(45,1) from (1,1)", polar{math.Pi / 4, sqrt2}, point{1, 1}, point{2, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := polar{tt.p.angle, tt.p.radius}
			if got := p.toCartesian(tt.start); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toCartesian() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_part2(t *testing.T) {
	rawLines := files.ReadAllLines("input.txt")
	type args struct {
		start        point
		rawLines     []string
		numAsteroids int
	}
	tests := []struct {
		name string
		args args
		want point
	}{
		{"2x1, n=1", args{point{0, 0}, []string{"##"}, 1}, point{1, 0}},
		{"3x1, n=2", args{point{0, 0}, []string{"###"}, 2}, point{2, 0}},
		{"3x2, n=3", args{point{0, 0}, []string{
			"#.#",
			"#.#",
		}, 3}, point{0, 1}},
		{"5x5, n=3", args{point{2, 2}, []string{
			"..#.#",
			"..#..",
			"..#..",
			"..#.#",
		}, 5}, point{2, 0}},
		{"input", args{point{27, 19}, rawLines, 200}, point{15, 13}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.args.start, tt.args.rawLines, tt.args.numAsteroids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
