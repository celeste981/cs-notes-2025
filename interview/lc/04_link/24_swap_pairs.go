package _4_link

func swapPairs(head *ListNode) *ListNode {
	dummy := &ListNode{
		Val:  0,
		Next: head,
	}
	cur := dummy
	for cur.Next != nil && cur.Next.Next != nil {
		a := cur.Next
		b := cur.Next.Next
		a.Next = b.Next
		b.Next = a
		cur.Next = b
		cur = a
	}
	return dummy.Next
	/*
		这个链表有四种情况
		第一种是空链表
		第二种是只有一个节点
		第三种是奇数个节点
		第四种是偶数个节点
	*/
}
