package _3_binary_tree

func levelOrder(root *TreeNode) [][]int {
	var ans [][]int
	if root == nil {
		return ans
	}
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		// 处理当前层
		size := len(queue)
		var level []int
		for i := 0; i < size; i++ {
			// 记录当前节点
			cur := queue[0]
			level = append(level, cur.Val)
			queue = queue[1:]
			// 将子结点入队
			if cur.Left != nil {
				queue = append(queue, cur.Left)
			}
			if cur.Right != nil {
				queue = append(queue, cur.Right)
			}
		}
		ans = append(ans, level)
	}
	return ans
}
