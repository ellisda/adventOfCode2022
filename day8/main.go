package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/samber/lo"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	grid := parseInput(f)

	// viewable := walkGrid(grid)
	viewable := walkGridFunc(grid)

	fmt.Println("viewable", viewable)

	// viewScore(grid, position{1, 2})
	viewScore(grid, position{3, 2})
	walkGridPart2(grid)
}

type position struct {
	row, col int
}

type mover func(*position)
type predicate func(grid [][]int, p position) bool

func walkGridFunc(grid [][]int) int {
	total := 0

	total += 2*len(grid) + 2*(len(grid[0])-2)
	fmt.Println("Edge Trees", total)

	seen := []position{}
	//Skip top and bottom row (the whole row is on the egde)
	for r := 1; r < len(grid)-1; r++ {

		fromLeft := moveAndCountTaller(func(p *position) { p.col++ }, position{r, 0}, grid)
		fromRight := moveAndCountTaller(func(p *position) { p.col-- }, position{r, len(grid[r]) - 1}, grid)

		seen = append(append(seen, fromLeft...), fromRight...)
	}

	for c := 1; c < len(grid[0])-1; c++ {
		fromTop := moveAndCountTaller(func(p *position) { p.row++ }, position{0, c}, grid)
		fromBottom := moveAndCountTaller(func(p *position) { p.row-- }, position{len(grid) - 1, c}, grid)

		seen = append(append(seen, fromTop...), fromBottom...)
	}
	unique := len(lo.Uniq(seen))
	total += unique

	fmt.Println("Inner Trees", unique, "Total", total)
	return total
}

// IsInner returns if the position is valid, and not an edge position
func (p position) IsInner(grid [][]int) bool {
	return p.row > 0 && p.row < len(grid)-1 &&
		p.col > 0 && p.col < len(grid[0])-1
}

func moveAndCountTaller(f mover, start position, grid [][]int) []position {
	tallerThanStart := []position{}
	maxYet := grid[start.row][start.col]

	p := start
	f(&p)

	for ; p.IsInner(grid); f(&p) {
		if grid[p.row][p.col] > maxYet {
			fmt.Println("New Taller Height", grid[p.row][p.col], "at", p, "PrevTallest", maxYet, "")
			maxYet = grid[p.row][p.col]
			tallerThanStart = append(tallerThanStart, p)
		}
	}
	return tallerThanStart
}

func walkGridPart2(grid [][]int) {
	bestViewScore := 0

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if candidate := viewScore(grid, position{i, j}); candidate > bestViewScore {
				bestViewScore = candidate
			}
		}
	}

	fmt.Println("Best yet:", bestViewScore)
}

func viewScore(grid [][]int, v position) int {
	right, left, up, down := 0, 0, 0, 0
	start := position{v.row, v.col + 1}
	for goingRight := start; goingRight.col < len(grid[v.row]); goingRight.col++ {
		self := grid[v.row][v.col]
		candidate := grid[goingRight.row][goingRight.col]
		if candidate < self {
			right++
		} else {
			right++
			break
		}
	}

	start = position{v.row, v.col - 1}
	for goingLeft := start; goingLeft.col >= 0; goingLeft.col-- {
		self := grid[v.row][v.col]
		candidate := grid[goingLeft.row][goingLeft.col]
		if candidate < self {
			left++
		} else {
			left++
			break
		}
	}

	start = position{v.row + 1, v.col}
	for goingDown := start; goingDown.row < len(grid[v.row]); goingDown.row++ {
		self := grid[v.row][v.col]
		candidate := grid[goingDown.row][goingDown.col]
		if candidate < self {
			down++
		} else {
			down++
			break
		}
	}

	start = position{v.row - 1, v.col}
	for goingUp := start; goingUp.row >= 0; goingUp.row-- {
		self := grid[v.row][v.col]
		candidate := grid[goingUp.row][goingUp.col]
		if candidate < self {
			up++
		} else {
			up++
			break
		}
	}
	total := up * down * left * right
	fmt.Println("ViewScore", total, v)
	return total
}

func filterMover(f mover, start position, check predicate, grid [][]int) []position {
	ret := []position{}

	p := start
	f(&p)

	for ; p.IsInner(grid); f(&p) {
		if check(grid, p) {
			fmt.Println("Matches Predicate", grid[p.row][p.col], "at", p)
			ret = append(ret, p)
		}
	}
	return ret
}

func walkGrid(grid [][]int) int {
	total := 0
	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if i == 0 || i == len(grid)-1 ||
				j == 0 || j == len(grid[i])-1 {
				fmt.Println("Edge Tree")
				total++
			} else if isTallestInAnyDir(grid, i, j) {
				fmt.Println("Viewable Tree")
				total++
			}
		}
	}
	return total
}

func isTallestInAnyDir(grid [][]int, i0, j0 int) bool {
	for i := i0 + 1; i < len(grid); i++ {
		if grid[i0][j0] > grid[i][j0] {
			if i == len(grid)-1 {
				return true
			}

		}
	}

	for i := i0 - 0; i >= 0 && grid[i0][j0] > grid[i][j0]; i++ {
		if i == 0 {
			return true
		}
	}

	for j := j0 + 1; j < len(grid[i0]) && grid[i0][j0] > grid[i0][j]; j++ {
		if j == len(grid[i0])-1 {
			return true
		}
	}

	for j := j0 - 1; j >= 0 && grid[i0][j0] > grid[i0][j]; j-- {
		if j == 0 {
			return true
		}
	}
	return false
}

func parseInput(f io.ReadSeekCloser) [][]int {

	ret := [][]int{}
	lines := []string{}
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		lines = append(lines, line)

		row := make([]int, len(line))
		for i, a := range line {
			row[i] = int(a) - 48
		}
		ret = append(ret, row)
	}
	return ret
}
