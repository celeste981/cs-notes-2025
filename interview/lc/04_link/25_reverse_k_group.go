package _4_link

func reverseKGroup(head *ListNode, k int) *ListNode {
	dummy := &ListNode{
		Val:  0,
		Next: head,
	}
	prev := dummy

	for true {
		// 每次找k个节点，找不到的时候才退出
		tail := prev
		for i := 0; i < k && tail != nil; i++ {
			tail = tail.Next
		}
		if tail == nil {
			break
		}
		next := tail.Next
		start := prev.Next

		// k 个一组反转链表
		cur := start
		prevNode := next
		for cur != next {
			tmp := cur.Next
			cur.Next = prevNode
			prevNode = cur
			cur = tmp
		}

		// 接回主链表
		prev.Next = tail
		prev = start
	}

	return dummy.Next
}
