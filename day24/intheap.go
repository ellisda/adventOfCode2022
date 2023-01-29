package main

// An IntHeap is a min-heap of ints.
type IntHeap struct {
	all  []move
	dest position
}

func (p position) DistanceFrom(target position) int {
	return abs(target.x-p.x) + abs(target.y-p.y)
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func (h IntHeap) Len() int { return len(h.all) }
func (h IntHeap) Less(i, j int) bool {
	//NOTE - If we don't prioritize distance, the BFS runs crazy slow and generates all sorts of walk-in-circles
	if h.all[i].stepsTo == h.all[j].stepsTo {
		return h.all[i].to.DistanceFrom(h.dest) < h.all[j].to.DistanceFrom(h.dest)
	}
	return h.all[i].stepsTo < h.all[j].stepsTo

	// di := h.all[i].to.DistanceFrom(h.dest)
	// dj := h.all[j].to.DistanceFrom(h.dest)
	// if di == dj {
	// 	//Tie breaker, prefer moves that have taken fewer steps to get there
	// 	return h.all[i].stepsTo < h.all[j].stepsTo
	// }
	// return di < dj

	//Try to blend distance and step count
	// di := h.all[i].to.DistanceFrom(h.dest) + h.all[i].stepsTo/10
	// dj := h.all[j].to.DistanceFrom(h.dest) + h.all[j].stepsTo/10
	// return di < dj
}

func (h IntHeap) Swap(i, j int) { h.all[i], h.all[j] = h.all[j], h.all[i] }

func (h *IntHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	h.all = append(h.all, x.(move))
}

func (h *IntHeap) Pop() any {
	old := h.all
	n := len(old)
	x := old[n-1]
	h.all = old[0 : n-1]
	return x
}
