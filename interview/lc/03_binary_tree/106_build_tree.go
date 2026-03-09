package _3_binary_tree

func buildTreePostOrder(inorder []int, postorder []int) *TreeNode {
	if len(postorder) == 0 {
		return nil
	}
	root := &TreeNode{
		Val:   postorder[len(postorder)-1],
		Left:  nil,
		Right: nil,
	}

	index := 0
	for inorder[index] != root.Val {
		index++
	}

	root.Left = buildTreePostOrder(inorder[:index], postorder[:index])
	root.Right = buildTreePostOrder(inorder[index+1:], postorder[index:len(postorder)-1])

	return root
}
