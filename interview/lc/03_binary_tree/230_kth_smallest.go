package _3_binary_tree

func kthSmallest(root *TreeNode, k int) int {
	cur := 0
	ans := 0
	dfsForKthSmallest(root, &cur, &ans, k)
	return ans
}

func dfsForKthSmallest(node *TreeNode, cur, ans *int, k int) {
	if node == nil {
		return
	}
	dfsForKthSmallest(node.Left, cur, ans, k)
	*cur = *cur + 1
	if *cur == k {
		*ans = node.Val
	}
	dfsForKthSmallest(node.Right, cur, ans, k)
}
