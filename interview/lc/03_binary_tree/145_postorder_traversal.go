package _3_binary_tree

func postorderTraversal(root *TreeNode) []int {
	var res []int
	postorder(root, &res)
	return res
}

func postorder(node *TreeNode, res *[]int) {
	if node == nil {
		return
	}
	postorder(node.Left, res)
	postorder(node.Right, res)
	*res = append(*res, node.Val)
}
