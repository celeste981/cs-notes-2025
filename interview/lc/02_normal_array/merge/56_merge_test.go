package merge

import (
	"reflect"
	"sort"
	"testing"
)

func TestMerge(t *testing.T) {
	tests := []struct {
		name string
		args [][]int
		want [][]int
	}{
		{"basic", [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}, [][]int{{1, 6}, {8, 10}, {15, 18}}},
		{"reorder", [][]int{{4, 7}, {1, 4}}, [][]int{{1, 7}}},
		{"no_interval", [][]int{{4, 7}, {1, 3}}, [][]int{{4, 7}, {1, 3}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Merge(tt.args)
			normalize(got)
			normalize(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Merge(%v) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func normalize(args [][]int) [][]int {
	for _, arg := range args {
		sort.Ints(arg)
	}
	sort.Slice(args, func(i, j int) bool {
		a, b := args[i], args[j]
		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] != b[k] {
				return a[k] < b[k]
			}
		}
		return len(a) < len(b)
	})
	return args
}
