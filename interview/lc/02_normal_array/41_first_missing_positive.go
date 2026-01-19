package _2_normal_array

func FirstMissingPositive(nums []int) int {
	// 因为nums.len == n，所以最多数n下，最坏情况是n+1，否则答案就是1..n

	// 给 1 .. n 的数放进数组里，其他大数不管
	// nums[i] == nums[nums[i] - 1]，这样的数是合格的
	for i := 0; i < len(nums); i++ {
		// 把 nums[i] 放到它应该在的位置 nums[i]-1
		for nums[i] >= 1 && nums[i] <= len(nums) && nums[nums[i]-1] != nums[i] {
			j := nums[i] - 1
			nums[i], nums[j] = nums[j], nums[i]
		}
	}

	// 遍历nums 数组，第一个不相符的数，就是答案，否则答案就是n+1
	for i := 0; i < len(nums); i++ {
		if nums[i] != i+1 {
			return i + 1
		}
	}

	return len(nums) + 1
}
