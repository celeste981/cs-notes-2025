package _8_random

func isSymmetricDup(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return dfsSymmetricDup(root.Left, root.Right)
}

func dfsSymmetricDup(left *TreeNode, right *TreeNode) bool {
	if left == nil && right == nil {
		return true
	}
	if left == nil || right == nil {
		return false
	}
	if left.Val != right.Val {
		return false
	}
	return dfsSymmetricDup(left.Left, right.Right) && dfsSymmetricDup(left.Right, right.Left)
}
