package _3_binary_tree

func invertTree(root *TreeNode) *TreeNode {
	if root == nil {
		return root
	}
	left := invertTree(root.Right)
	right := invertTree(root.Left)
	root.Left = left
	root.Right = right
	return root
}
