package _4_link

func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	cur := 0
	step := 0
	dummy := &ListNode{}
	head := dummy
	for l1 != nil || l2 != nil || step != 0 {
		cur = step
		if l1 != nil {
			cur += l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			cur += l2.Val
			l2 = l2.Next
		}
		step = cur / 10
		cur %= 10
		head.Next = &ListNode{
			Val:  cur,
			Next: nil,
		}
		head = head.Next
	}

	return dummy.Next
}
