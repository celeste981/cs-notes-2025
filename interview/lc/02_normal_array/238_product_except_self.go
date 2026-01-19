package _2_normal_array

func ProductExceptSelf(nums []int) []int {
	// 当前元素的结果的等于 左乘积 * 右乘积
	ans := make([]int, len(nums))
	// 先算左乘积
	ans[0] = 1
	for i := 1; i < len(nums); i++ {
		ans[i] = ans[i-1] * nums[i-1]
	}

	// 再算右乘积
	cur := 1
	for i := len(nums) - 2; i >= 0; i-- {
		cur *= nums[i+1]
		ans[i] = ans[i] * cur
	}

	return ans
}
