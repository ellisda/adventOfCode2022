package linkedlist

import (
	"github.com/samber/lo"
)

func NewCircular(s []int) *Node {
	numNodes := len(s)
	if numNodes == 0 {
		return nil
	}
	// c := make(circular, len(s))
	ret := &Node{s[0], nil, &numNodes}
	prev := ret
	for i := 1; i < len(s); i++ {
		prev.Next = &Node{s[i], nil, &numNodes}
		prev = prev.Next
	}
	prev.Next = ret //circular
	return ret
}

type Node struct {
	Value  int
	Next   *Node
	Length *int
	// origIndex int //use a secondary structure to remember inital order
}

// func (c circular) Copy() *circular {
// 	ret := make(circular, len(c))
// 	if n := copy(ret, c); n != len(c) {
// 		log.Fatal("copy failure")
// 	}

// 	return &ret
// }

func (n Node) Len() int {
	return *n.Length
}

func (n *Node) DeepCopy() *Node {
	return NewCircular(n.Ints())
}

func (n *Node) Order() []*Node {
	ret := make([]*Node, *n.Length)
	cur := n
	for i := 0; i < len(ret); i++ {
		ret[i] = cur
		cur = cur.Next
	}
	return ret
}

func (n *Node) Ints() []int {
	ret := []int{n.Value}

	c := n.Next
	for c != n {
		ret = append(ret, c.Value)
		c = c.Next
	}
	return ret
}

func (n *Node) Distr() map[int]int {
	return lo.CountValues(n.Ints())
}

// CycleOnce returns a decrypted sequence
func (n *Node) CycleN(numTimes int) {
	order := n.Order()

	for i := 0; i < numTimes; i++ {
		for _, c := range order {
			c.MoveItem(0)
		}
	}
}

func (n *Node) FindItemsFromZero(indices ...int) []int {
	z := n
	for z.Value != 0 {
		z = z.Next
	}
	ret := make([]int, len(indices))
	for i, idx := range indices {
		ret[i] = z.FindItem(idx).Value
	}
	return ret
}

func (n *Node) FindItem(idx_orig int) *Node {
	idx := idx_orig % (*(n.Length)) //no need to walk 2+ laps
	// fmt.Println("FindItem", idx_orig, idx)
	if idx < 0 {
		idx = *n.Length + idx //walk forward till we get almost around
	}

	seen, cur := 0, n
	for seen < idx {
		cur = cur.Next
		seen++
	}
	return cur
}

func (n *Node) MoveItem(i int) {
	// if i < 0 {
	// 	panic("need handle negative Find?")
	// } else if i == 0 {
	// 	panic("need handle Move root?")
	// }

	left := n.FindItem(i - 1)
	target := left.Next

	if target.Value%((*n.Length)-1) == 0 {
		return
	}

	//Seek the one to the left of target position (for insertion)
	seekFromTarget := target.Value % (*(n.Length) - 1)
	if seekFromTarget < 0 {
		seekFromTarget--
	}

	// if seekFromTarget%((*n.Length)-1) == 0 {
	// 	return
	// }

	right := target.FindItem(seekFromTarget)

	//Handle special case where we move 1
	if target == right {
		right = right.Next
	}

	//Remove target from prev position
	left.Next = left.Next.Next

	//Insert new node in new position
	target.Next = right.Next
	right.Next = target
}
