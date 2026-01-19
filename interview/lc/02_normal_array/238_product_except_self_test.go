package _2_normal_array

import (
	"reflect"
	"testing"
)

func TestProductExceptSelf(t *testing.T) {
	type args struct {
		nums []int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{"basic", struct{ nums []int }{nums: []int{1, 2, 3, 4}}, []int{24, 12, 8, 6}},
		{"contains_zero", struct{ nums []int }{nums: []int{-1, 1, 0, -3, 3}}, []int{0, 0, 9, 0, 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProductExceptSelf(tt.args.nums); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProductExceptSelf() = %v, want %v", got, tt.want)
			}
		})
	}
}
