package _3_binary_tree

func isBalanced(root *TreeNode) bool {
	return dfsForIsBalanced(root) != -1
}

func dfsForIsBalanced(node *TreeNode) int {
	if node == nil {
		return 0
	}
	left := dfsForIsBalanced(node.Left)
	// 提前短路，不需要维护全局变量
	if left == -1 {
		return -1
	}
	right := dfsForIsBalanced(node.Right)
	if right == -1 {
		return -1
	}
	if left-right > 1 || right-left > 1 {
		return -1
	}
	return max(left, right) + 1
}
