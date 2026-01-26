package _3_binary_tree

func isSymmetric(root *TreeNode) bool {
	if root == nil {
		return true
	}
	return dfsForIsSymmetric(root.Left, root.Right)
}

func dfsForIsSymmetric(left *TreeNode, right *TreeNode) bool {
	if left == nil && right == nil {
		return true
	}
	if (left == nil && right != nil) || (left != nil && right == nil) {
		return false
	}
	if left.Val != right.Val {
		return false
	}
	return dfsForIsSymmetric(left.Left, right.Right) && dfsForIsSymmetric(left.Right, right.Left)
}
