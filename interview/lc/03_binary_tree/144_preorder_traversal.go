package _3_binary_tree

func preorderTraversal(root *TreeNode) []int {
	var res []int
	preorder(root, &res)
	return res
}

func preorder(node *TreeNode, res *[]int) {
	if node == nil {
		return
	}
	*res = append(*res, node.Val)
	preorder(node.Left, res)
	preorder(node.Right, res)
}
