package _6_trick

func singleNumber(nums []int) int {
	var ans int
	for num := range nums {
		ans ^= nums[num]
	}
	return ans
}
