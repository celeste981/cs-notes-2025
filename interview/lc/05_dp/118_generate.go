package _5_dp

func generate(numRows int) [][]int {
	var ans [][]int
	ans = append(ans, []int{1})
	if numRows == 1 {
		return ans
	}
	for i := 1; i < numRows; i++ {
		var level []int
		level = append(level, 1)
		for j := 1; j < len(ans[i-1]); j++ {
			level = append(level, ans[i-1][j-1]+ans[i-1][j])
		}
		level = append(level, 1)
		ans = append(ans, level)
	}
	return ans
}
