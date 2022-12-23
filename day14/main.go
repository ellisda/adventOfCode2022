package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type grid [][]bool

type point struct {
	x int
	y int
}

func main() {
	//begin bernard parse
	file, _ := os.Open("input.txt")
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var input []string

	for scanner.Scan() {
		v := scanner.Text()
		input = append(input, v)
	}
	//end
	multis := parseInput(input)
	fmt.Println(multis)
	g := buildGrid(multis)
	before := g.Count()
	fmt.Println("Before Sand", before)
	g.FillWithSand(point{500, 0})
	after := g.Count()
	fmt.Println("Added sand", after-before)

	g2 := buildGrid(multis)
	g2.AddPart2Floor()
	b2 := g2.Count()
	fmt.Println("Before Sand", b2)
	g2.FillWithSandPart2(point{500, 0})
	a2 := g2.Count()
	fmt.Println("Part2 Added sand", a2-b2)

}

func buildGrid(multis [][]point) grid {
	max := 700
	g := make(grid, max)
	for i := 0; i < max; i++ {
		g[i] = make([]bool, max)
	}

	for _, m := range multis {
		for i := 0; i+1 < len(m); i++ {
			fmt.Print("Segment")
			for _, p := range m[i].PointsToward(m[i+1]) {
				g[p.x][p.y] = true
				fmt.Print(p)
			}
			fmt.Println()
		}
	}

	return g
}

func (g grid) Count() int {
	ret := 0
	for x := 0; x < len(g); x++ {
		for y := 0; y < len(g[0]); y++ {
			if g[x][y] {
				ret++
			}
		}
	}
	return ret
}

func (g grid) MaxY() int {
	max := 0
	for x := 0; x < len(g); x++ {
		for y := 0; y < len(g[0]); y++ {
			if g[x][y] && y > max {
				max = y
			}
		}
	}
	return max
}

func (g grid) FillWithSand(source point) {
	// add sand from source

	cur := source
	// stop when sand falls past y > 500
	for cur.y < 598 {
		// follow rules: try to move down,
		switch {
		// can move down
		case g[cur.x][cur.y+1] == false:
			cur.y = cur.y + 1

		// can move diag down left
		case g[cur.x-1][cur.y+1] == false:
			cur.x--
			cur.y++

		// can move diag down right
		case g[cur.x+1][cur.y+1] == false:
			cur.x++
			cur.y++

		default:
			//sand hardens and fills in place in grid
			g[cur.x][cur.y] = true
			cur = source
		}
	}
}

func (g grid) AddPart2Floor() {
	floor := g.MaxY() + 2
	fmt.Println("Adding floor at", floor)
	for x := 0; x < len(g); x++ {
		g[x][floor] = true
	}
}

func (g grid) FillWithSandPart2(source point) {
	// add sand from source

	cur := source
	// stop when sand falls past y > 500
	for {
		if cur.x == 599 {
			// fmt.Println("huh")
		}
		// follow rules: try to move down,
		switch {

		// can move down
		case g[cur.x][cur.y+1] == false:
			cur.y = cur.y + 1

		// can move diag down left
		case g[cur.x-1][cur.y+1] == false:
			cur.x--
			cur.y++

		// can move diag down right
		case g[cur.x+1][cur.y+1] == false:
			cur.x++
			cur.y++

		default:
			//sand hardens and fills in place in grid
			g[cur.x][cur.y] = true
			if cur == source {
				return
			}
			cur = source
		}
	}
}

func (p point) PointsToward(dest point) []point {
	ret := []point{p}
	for ret[len(ret)-1] != dest {
		ret = append(ret, ret[len(ret)-1].StepToward(dest))
	}
	return ret
}

func (p point) Less(p2 point) bool {
	return p.x < p2.x || p.y < p2.y
}

// StepToward returns a new point that is a unit step towards dest
func (p point) StepToward(dest point) point {
	switch {

	case dest.x > p.x:
		return point{p.x + 1, p.y}
	case dest.x < p.x:
		return point{p.x - 1, p.y}

	case dest.y > p.y:
		return point{p.x, p.y + 1}
	case dest.y < p.y:
		return point{p.x, p.y - 1}

	default:
		panic("bad")
	}
}

func parseInput(lines []string) [][]point {

	multilines := [][]point{}
	for _, line := range lines {
		multi := []point{}
		for _, pair := range strings.Split(line, " -> ") {
			p := strings.Split(pair, ",")
			x, _ := strconv.Atoi(p[0])
			y, _ := strconv.Atoi(p[1])
			multi = append(multi, point{x, y})
		}

		multilines = append(multilines, multi)
	}
	return multilines
}
