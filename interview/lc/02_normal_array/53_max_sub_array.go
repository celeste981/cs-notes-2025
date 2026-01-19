package _2_normal_array

func MaxSubArray(nums []int) int {
	var ans, cur = nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		num := nums[i]
		if num > cur+num {
			cur = num
		} else {
			cur += num
		}
		if cur > ans {
			ans = cur
		}
	}
	return ans
}
