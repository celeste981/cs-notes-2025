package _1_array_hash

func LongestConsecutive(nums []int) int {
	var res = 0

	// 用 set 来存储数字是否出现
	set := make(map[int]struct{})

	for _, num := range nums {
		set[num] = struct{}{}
	}

	for num := range set {
		if _, ok := set[num-1]; !ok {
			len := 0
			for {
				if _, ok := set[num]; !ok {
					break
				}
				len++
				num++
			}
			if len > res {
				res = len
			}
		}
	}

	return res
}
