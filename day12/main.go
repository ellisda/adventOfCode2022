package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
)

const (
	Start rune = 'S'
	End   rune = 'E'
)

type square struct {
	x, y      int
	elevation rune
	backTrack *square
	// backScore func(prev *square) int //number of moves it took to get here from backTrack
	visited bool
}

type grid [][]*square

type move struct {
	from *square
	to   *square
}

type moves []move

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	grid, start, end := parseInput(f)

	grid.PrintBoard()

	fmt.Println("start", start)
	fmt.Println("end", end)
	// grid.ProcessSquare(start)

	q := &moves{}
	grid.WalkAndEnqueueBackwards(end, q)
	grid.DequeueAndBuildBacktrack(q)

	fmt.Println("End", end, end.BackScore())
	fmt.Println(countBackTrack(start))

	grid.PrintBacktrack(start)

}

func (ms *moves) Enqueue(moves ...move) {
	*ms = append(*ms, moves...)
}
func (ms *moves) Dequeue() (move, bool) {
	if len(*ms) == 0 {
		return move{}, false
	}

	m := (*ms)[0]
	*ms = (*ms)[1:]
	return m, true
}

func (g grid) PrintBoard() {
	fmt.Println("Grid:")
	for _, row := range g {
		for _, s := range row {
			fmt.Print(string(s.elevation))
		}
		fmt.Println()
	}
}

func (s *square) BackScore() int {
	return s.backScore(false)
}

func (s *square) backScore(debug bool) int {
	switch {
	case s.elevation == End:
		return 0
	case s.backTrack == nil:
		return math.MaxInt
	default:
		if debug {
			fmt.Print("score_prev", s.backTrack)
		}
		return 1 + s.backTrack.backScore(debug)
	}
}

func (g grid) WalkAndEnqueueBackwards(here *square, q *moves) {
	if here.elevation == Start || here.visited {
		return
	}
	for _, candidate := range g.getCandidates(here) {
		if canDescend(here, candidate) {
			q.Enqueue(move{from: here, to: candidate})
		}
	}

}

func (g grid) getCandidates(here *square) []*square {
	candidates := []*square{}
	if here.x-1 >= 0 {
		candidates = append(candidates, g[here.y][here.x-1])
	}
	if here.x+1 < len(g[here.y]) {
		candidates = append(candidates, g[here.y][here.x+1])
	}
	if here.y+1 < len(g) {
		candidates = append(candidates, g[here.y+1][here.x])
	}
	if here.y-1 >= 0 {
		candidates = append(candidates, g[here.y-1][here.x])
	}
	return candidates
}

func canDescend(here, candidate *square) bool {
	switch {
	case here.elevation == End:
		return candidate.elevation == rune('z')

	case candidate.elevation == Start:
		return true

	default:
		// can only descend one level without gear
		return here.elevation <= candidate.elevation+1
	}
}

// func (g grid) WalkAndEnqueue(here *square, q *moves) {
// 	if here.elevation == End || here.visited {
// 		return
// 	}
// 	for _, candidate := range g.getCandidates(here) {
// 		q.Enqueue(move{from: here, to: candidate})
// 	}

// }

func (g grid) DequeueAndBuildBacktrack(q *moves) {
	for m, ok := q.Dequeue(); ok; m, ok = q.Dequeue() {
		here, candidate := m.from, m.to

		// if here.x == 2 && here.y == 2 && candidate.x == 2 && candidate.y == 3 {
		// 	fmt.Println("???  ", candidate, "from", here)
		// }

		if here.x == 2 && here.y == 1 && candidate.x == 2 && candidate.y == 2 {
			fmt.Println("???  ", candidate, "from", here)
		}

		// if candidate.x == 2 && candidate.y == 3 {
		// 	fmt.Println("crit move (before)", candidate.BackScore(), candidate, here)
		// 	if here.x == 2 && here.y == 2 {
		// 		fmt.Println("MAGIC")
		// 	}
		// }
		if updateBackTrack(here, candidate) {
			if !candidate.visited {
				g.WalkAndEnqueueBackwards(candidate, q)
			}
			// g.updateBackScore(candidate)
		}
		// if candidate.x == 2 && candidate.y == 3 {
		// 	fmt.Println("crit move (after)", candidate.BackScore(), candidate, here)
		// }

		here.visited = true
	}
}

func updateBackTrack(here, candidate *square) bool {

	// if here.x == 2 && here.y == 2 && candidate.x == 2 && candidate.y == 3 {
	// if here.x == 2 && here.y == 1 && candidate.x == 2 && candidate.y == 2 {
	// if here.x == 2 && here.y == 0 && candidate.x == 2 && candidate.y == 1 {
	// if //(here.x == 1 && here.y == 0 && candidate.x == 2 && candidate.y == 0) ||
	// here.x == 0 && here.y == 0 && candidate.x == 1 && candidate.y == 0 {
	// 	fmt.Println("Checking for Better Backtrack for ", candidate, "from", here)
	// 	fmt.Println("crit check")
	// 	fmt.Print("here Score: ", here)
	// 	h := here.backScore(true)
	// 	fmt.Println(" --- FINAL", h)

	// 	fmt.Print("\ncanidate Score: ", candidate)
	// 	c := candidate.backScore(true)
	// 	fmt.Println(" --- FINAL", c)
	// }

	hScore := here.BackScore()
	cScore := candidate.BackScore()
	if hScore < cScore-1 {
		//Found new best path to reach candidate
		candidate.backTrack = here
		// candidate.backScore = ScorePlusOne(here)

		fmt.Println("New Best Backtrack for ", candidate, "from", here)
		return true
	}
	return false
}

func canHikeUp(here, candidate *square) bool {
	switch {
	case here.elevation == Start:
		return true
	case candidate.elevation == End:
		return here.elevation == rune('z')
	default:
		return here.elevation+1 >= candidate.elevation
	}
}

// func (g grid) updateBackScore(here *square) {
// 	if here.elevation == End {
// 		return
// 	}

// 	for _, c := range g.getCandidates(here) {
// 		if c.backTrack == here {
// 			c.backScore = here.backScore + 1
// 			fmt.Println("Updating BackScore", c)
// 			g.updateBackScore(c)
// 		}
// 	}
// }

func countBackTrack(end *square) int {
	total := 0
	for next := end.backTrack; next != nil; {
		total++
		fmt.Println("counting", string(next.elevation), total, next)
		next = next.backTrack
	}
	return total
}

func (g grid) PrintBacktrack(end *square) {
	fmt.Println("Back Track:")
	for _, row := range g {
		for _, s := range row {
			fwd := g.GetForwardTrack(s)
			switch {
			// case s.elevation == Start:
			// 	fmt.Print("S")
			case s.elevation == End:
				fmt.Print("E")
			case g.IsOnPathToEnd(s, end) == false:
				fmt.Print(".")
			case fwd == nil:
				fmt.Print(" ")
			case fwd.x-s.x == 1:
				fmt.Print(">")
			case fwd.x-s.x == -1:
				fmt.Print("<")

			case fwd.y-s.y == 1:
				fmt.Print("v")
			case fwd.y-s.y == -1:
				fmt.Print("^")
			default:
				panic("??")
			}
		}
		fmt.Println()
	}
}

func (g grid) IsOnPathToEnd(s *square, end *square) bool {
	for x := end.backTrack; x != nil; x = x.backTrack {
		if x == s {
			return true
		}
	}
	return false
}

func (g grid) GetForwardTrack(here *square) *square {
	for _, row := range g {
		for _, s := range row {
			if s.backTrack == here {
				return s
			}
		}
	}
	return nil
}

func parseInput(f io.ReadSeekCloser) (ret grid, start *square, end *square) {
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	y := 0
	for s.Scan() {
		line := s.Text()
		row := make([]*square, len(line))

		for x, c := range line {
			row[x] = &square{x, y, c, nil, false}
			switch c {
			case Start:
				start = row[x]
			case End:
				end = row[x]
			}
		}
		ret = append(ret, row)
		y++
	}

	return
}
