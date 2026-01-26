package _3_binary_tree

import (
	"math"
)

func isValidBST(root *TreeNode) bool {
	return dfsForIsValidBST(root, math.MinInt, math.MaxInt)
}

func dfsForIsValidBST(node *TreeNode, min, max int) bool {
	if node == nil {
		return true
	}
	if node.Val >= max || node.Val <= min {
		return false
	}
	return dfsForIsValidBST(node.Left, min, node.Val) && dfsForIsValidBST(node.Right, node.Val, max)
}
