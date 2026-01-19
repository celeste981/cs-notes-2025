package _2_normal_array

func Rotate(nums []int, k int) []int {
	if k == 0 {
		return nums
	}
	ans := make([]int, len(nums))
	k %= len(nums)

	for i := 0; i < k; i++ {
		ans[i] = nums[len(nums)-k+i]
	}

	for i := k; i < len(nums); i++ {
		ans[i] = nums[i-k]
	}

	nums = ans
	return nums
}
