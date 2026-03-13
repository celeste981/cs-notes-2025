package _4_link

func getIntersectionNode(headA, headB *ListNode) *ListNode {
	if headA == nil {
		return headB
	}
	if headB == nil {
		return headA
	}
	pA := headA
	pB := headB
	for pA != pB {
		if pA == nil {
			pA = headB
		} else {
			pA = pA.Next
		}

		if pB == nil {
			pB = headA
		} else {
			pB = pB.Next
		}
	}
	return pA
}
