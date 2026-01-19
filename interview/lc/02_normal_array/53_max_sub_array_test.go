package _2_normal_array

import (
	"testing"
)

func TestMaxSubArray(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		want int
	}{
		{"basic", []int{-2, 1, -3, 4, -1, 2, 1, -5, 4}, 6},
		{"one_num", []int{1}, 1},
		{"all_num", []int{5, 4, -1, 7, 8}, 23},
		{"all_neg", []int{-5, -4, -1}, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaxSubArray(tt.nums)
			if got != tt.want {
				t.Errorf("MaxSubArray(%v) = %v, want %v", tt.nums, got, tt.want)
			}
		})
	}
}
