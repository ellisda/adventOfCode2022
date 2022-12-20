package main

import (
	"container/heap"
	"fmt"
	"log"
)

type position struct {
	x int
	y int
}

func (s square) position() position {
	return position{s.x, s.y}
}

func (g grid) Shortest(start, end *square) int {
	visited := make(map[position]int)
	visited[start.position()] = 0

	directions := make(map[position]string)

	neighbors := g.getCandidates(start)
	// Create a priority queue, put the items in it, and
	// establish the priority queue (heap) invariants.
	pq := make(PriorityQueue, 0)
	for i, dest := range neighbors {
		if canHikeUp(start, dest) {
			pq = append(pq, &Item{
				value:    move{from: start, to: dest},
				priority: 1,
				index:    i,
			})
		}
	}
	heap.Init(&pq)

	var next *Item
	for len(pq) > 0 {
		next = heap.Pop(&pq).(*Item)

		if best, ok := visited[next.value.to.position()]; ok {
			if next.priority < best {
				log.Fatal("Found even better")
			}
			continue
		} else {
			visited[next.value.to.position()] = next.priority
			directions[next.value.to.position()] = direction(next.value)
			if next.value.to == end {
				break
			}
		}
		// fmt.Println("Visited:", next.priority, next.value.to)
		if len(pq)%100000 == 0 {
			fmt.Println("Len(queue)", len(pq))
		}

		for _, c := range g.getCandidates(next.value.to) {
			if canHikeUp(next.value.to, c) {
				if _, beenThere := visited[c.position()]; beenThere || next.value.from == c {
					continue //don't re-enqueue places we've just come from or places we've sealed
				}
				heap.Push(&pq, &Item{
					priority: next.priority + 1,
					value:    move{from: next.value.to, to: c},
				})
			}
		}
	}

	printVisited(directions, visited, end)

	return visited[end.position()]
}

func direction(m move) string {
	switch {
	case m.from.x < m.to.x:
		return ">"
	case m.from.y < m.to.y:
		return "v"
	case m.from.x > m.to.x:
		return "<"
	case m.from.y > m.to.y:
		return "^"
	default:
		panic("bad choices")
	}
}

func printVisited(directions map[position]string, v map[position]int, end *square) {
	// sorted := lo.PickByValues(v, []int{352, 351})

	// for k, v := range sorted {
	// 	fmt.Println(v, k)
	// }

	delta := 5

	for y := end.y - delta; y < end.y+delta; y++ {
		for x := end.x - delta; x < end.x+delta; x++ {
			fmt.Print(x, y, " *", v[position{x, y}], "*, ")
		}
		fmt.Println()
	}

	fmt.Println("Direcions")
	for y := end.y - delta; y < end.y+delta; y++ {
		for x := end.x - delta; x < end.x+delta; x++ {
			if s, ok := directions[position{x, y}]; ok {
				fmt.Print(s)
			} else {
				fmt.Print(".")
			}
		}
		if y == end.y {
			fmt.Print(" END position at ", end.x, end.y)
		}
		fmt.Println()
	}

}

// func (g grid) GetNextMoves(prev move) []move {

// }

// An Item is something we manage in a priority queue.
type Item struct {
	value    move // The value of the item; arbitrary.
	priority int  // The priority of the item in the queue.
	// The index is needed by update and is maintained by the heap.Interface methods.
	index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority < pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value move, priority int) {
	item.value = value
	item.priority = priority
	heap.Fix(pq, item.index)
}
