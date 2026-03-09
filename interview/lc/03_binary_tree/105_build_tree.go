package _3_binary_tree

func buildTree(preorder []int, inorder []int) *TreeNode {
	if len(preorder) == 0 {
		return nil
	}
	root := &TreeNode{
		Val:   preorder[0],
		Left:  nil,
		Right: nil,
	}

	// 找出左子树和右子树

	index := 0
	for inorder[index] != root.Val {
		index++
	}
	root.Left = buildTree(preorder[1:index+1], inorder[0:index])
	root.Right = buildTree(preorder[index+1:], inorder[index+1:])
	return root
}
