package _4_link

func reverseList(head *ListNode) *ListNode {
	if head == nil {
		return nil
	}
	cur := head
	var prev *ListNode
	for cur != nil {
		next := cur.Next
		cur.Next = prev
		prev = cur
		cur = next
	}
	return prev
}
