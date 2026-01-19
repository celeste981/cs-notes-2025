package _2_normal_array

import (
	"sort"
)

func Merge(intervals [][]int) [][]int {
	// 按左边数字来排序
	//for _, interval := range intervals {
	//	sort.Ints(interval)
	//}
	sort.Slice(intervals, func(i, j int) bool {
		a, b := intervals[i], intervals[j]
		return a[0] < b[0]
	})

	ans := make([][]int, 0)
	cur := intervals[0]
	for _, interval := range intervals {
		if interval[0] <= cur[1] {
			if interval[1] > cur[1] {
				cur[1] = interval[1]
			}
		} else {
			ans = append(ans, cur)
			cur = interval
		}
	}
	ans = append(ans, cur)

	return ans
}
