package _6_trick

func majorityElement(nums []int) int {
	// 投票法配对
	var count, candidate int
	for _, num := range nums {
		if count == 0 {
			candidate = num
		}
		if num == candidate {
			count++
		} else {
			count--
		}
	}
	return candidate
}
