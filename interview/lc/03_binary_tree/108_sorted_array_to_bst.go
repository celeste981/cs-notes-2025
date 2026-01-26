package _3_binary_tree

func sortedArrayToBST(nums []int) *TreeNode {
	//用有序数组的中点建根节点，用左半数组递归建左子树，用右半数组递归建右子树
	return dfsForSortedArrayToBST(nums, 0, len(nums)-1)
}

func dfsForSortedArrayToBST(nums []int, l, r int) *TreeNode {
	if r < l {
		return nil
	}
	mid := l + (r-l)/2
	left := dfsForSortedArrayToBST(nums, l, mid-1)
	right := dfsForSortedArrayToBST(nums, mid+1, r)
	return &TreeNode{
		Val:   nums[mid],
		Left:  left,
		Right: right,
	}
}
