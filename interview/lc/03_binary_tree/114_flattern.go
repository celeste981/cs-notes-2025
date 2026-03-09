package _3_binary_tree

//func flatten(root *TreeNode) {
//	if root == nil {
//		return
//	}
//	var nodes []int
//	// dfs
//	dfsForFlatten(root, &nodes)
//	nodes = nodes[1:]
//	cur := root
//	for _, node := range nodes {
//		cur.Left = nil
//		cur.Right = &TreeNode{
//			Val:   node,
//			Left:  nil,
//			Right: nil,
//		}
//		cur = cur.Right
//	}
//}
//
//func dfsForFlatten(node *TreeNode, nodes *[]int) {
//	if node == nil {
//		return
//	}
//	*nodes = append(*nodes, node.Val)
//	dfsForFlatten(node.Left, nodes)
//	dfsForFlatten(node.Right, nodes)
//}

func flatten(root *TreeNode) {
	if root == nil {
		return
	}
	left := root.Left
	right := root.Right
	flatten(left)
	flatten(right)
	if left == nil {
		return
	}
	root.Left = nil
	root.Right = left
	cur := left
	for cur.Right != nil {
		cur = cur.Right
	}
	cur.Left = nil
	cur.Right = right
}
