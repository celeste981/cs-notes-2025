package _3_binary_tree

func rightSideView(root *TreeNode) []int {
	var ans []int
	if root == nil {
		return ans
	}
	var level [][]int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		var curLevel []int
		size := len(queue)
		for i := 0; i < size; i++ {
			cur := queue[0]
			curLevel = append(curLevel, cur.Val)
			queue = queue[1:]
			if cur.Left != nil {
				queue = append(queue, cur.Left)
			}
			if cur.Right != nil {
				queue = append(queue, cur.Right)
			}
		}
		level = append(level, curLevel)
	}
	for i := 0; i < len(level); i++ {
		ans = append(ans, level[i][len(level[i])-1])
	}
	return ans
}
