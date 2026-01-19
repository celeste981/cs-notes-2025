package _1_array_hash

import (
	"reflect"
	"sort"
	"testing"
)

func TestGroupAnagrams(t *testing.T) {
	tests := []struct {
		name string
		strs []string
		want [][]string
	}{
		{"basic", []string{"eat", "tea", "tan", "ate", "nat", "bat"}, [][]string{{"bat"}, {"nat", "tan"}, {"ate", "eat", "tea"}}}} //{"no result", []int{1, 2, 3}, 9, nil},

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GroupAnagrams(tt.strs)
			normalize(got)
			normalize(tt.want)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GroupAnagrams(%v) = %v, want %v", tt.strs, got, tt.want)
			}
		})
	}

}

func normalize(res [][]string) {
	for _, v := range res {
		sort.Strings(v)
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i][0] > res[j][0]
	})
}
