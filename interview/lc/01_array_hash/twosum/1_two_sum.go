package twosum

func TwoSum(nums []int, target int) []int {
	var res []int
	for i, num := range nums {
		for j := i + 1; j < len(nums); j++ {
			if target == num+nums[j] {
				res = []int{i, j}
				break
			}
		}
	}
	return res
}
