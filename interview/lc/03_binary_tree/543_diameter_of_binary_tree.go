package _3_binary_tree

func diameterOfBinaryTree(root *TreeNode) int {
	var ans int
	dfsForDiameterOfBinaryTree(root, &ans)
	return ans
}

func dfsForDiameterOfBinaryTree(node *TreeNode, ans *int) int {
	if node == nil {
		return 0
	}
	left := dfsForDiameterOfBinaryTree(node.Left, ans)
	right := dfsForDiameterOfBinaryTree(node.Right, ans)
	*ans = max(*ans, left+right)
	return max(left, right) + 1
}
