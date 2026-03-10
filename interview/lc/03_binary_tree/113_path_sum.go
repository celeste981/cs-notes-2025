package _3_binary_tree

func pathSum(root *TreeNode, targetSum int) [][]int {
	var ans [][]int
	var path []int
	dfsForPathSum(root, targetSum, &ans, &path)
	return ans
}

func dfsForPathSum(node *TreeNode, targetSum int, ans *[][]int, path *[]int) {
	if node == nil {
		return
	}
	*path = append(*path, node.Val)
	targetSum -= node.Val

	if node.Left == nil && node.Right == nil && 0 == targetSum {
		tmp := make([]int, len(*path))
		copy(tmp, *path)
		*ans = append(*ans, tmp)
	}

	dfsForPathSum(node.Left, targetSum, ans, path)
	dfsForPathSum(node.Right, targetSum, ans, path)

	*path = (*path)[:len(*path)-1]
}
