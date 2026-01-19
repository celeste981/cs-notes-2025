package _2_normal_array

import (
	"reflect"
	"testing"
)

func TestRotate(t *testing.T) {
	tests := []struct {
		name string
		nums []int
		k    int
		want []int
	}{
		{"basic", []int{1, 2, 3, 4, 5, 6, 7}, 3, []int{5, 6, 7, 1, 2, 3, 4}},
		{"basic", []int{-1, -100, 3, 99}, 2, []int{3, 99, -1, -100}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Rotate(tt.nums, tt.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Rotate(%v, %v) = %v, want %v", tt.nums, tt.k, got, tt.want)
			}
		})
	}
}
