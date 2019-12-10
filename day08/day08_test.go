package main

import (
	"reflect"
	"testing"
)

func Test_part2(t *testing.T) {
	type args struct {
		width  int
		height int
		ints   []int
	}
	tests := []struct {
		name string
		args args
		want [][]int
	}{
		{"1x1: 1 layer: w", args{1, 1, []int{White}}, [][]int{{White}}},
		{"1x1: 1 layer: b", args{1, 1, []int{Black}}, [][]int{{Black}}},
		{"1x1: 2 layers: bw", args{1, 1, []int{Black, White}}, [][]int{{Black}}},
		{"2x2: 4 layers: bw", args{2, 2,
			[]int{
				0, 2, 2, 2,
				1, 1, 2, 2,
				2, 2, 1, 2,
				0, 0, 0, 0}},
			[][]int{{Black, White}, {White, Black}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := part2(tt.args.ints, tt.args.width, tt.args.height); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("part2() = %v, want %v", got, tt.want)
			}
		})
	}
}
