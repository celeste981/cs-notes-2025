package _5_dp

func rob(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	if len(nums) == 1 {
		return nums[0]
	}
	var p1, p2 int
	for _, num := range nums {
		cur := max(p1, p2+num)
		p2 = p1
		p1 = cur
	}
	return p1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
