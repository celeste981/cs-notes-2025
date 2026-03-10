package _3_binary_tree

func pathSum437(root *TreeNode, targetSum int) int {
	if root == nil {
		return 0
	}
	return countSum(root, targetSum) + pathSum437(root.Left, targetSum) + pathSum437(root.Right, targetSum)
}

func countSum(node *TreeNode, target int) int {
	if node == nil {
		return 0
	}
	ans := 0
	target -= node.Val
	if target == 0 {
		ans++
	}
	ans += countSum(node.Left, target)
	ans += countSum(node.Right, target)
	return ans
}
