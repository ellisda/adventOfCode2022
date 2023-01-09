package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

const (
	N = iota
	S
	W
	E
)

type pos struct {
	x int
	y int
}
type elf struct {
	loc pos
	// steps int
}

type gang []elf

var POISON *elf = &elf{}

func main() {
	in := "input.txt"
	file, _ := os.Open(in)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// fmt.Println(strings.Join(lines, "\n"))

	g := ParseElves(lines)
	g10 := g.MoveN(10)
	fmt.Println("Part1", g10.CountSparseEmpties())

	fmt.Println("Part2")
	g.MoveN(1200)
}

func ParseElves(lines []string) gang {
	ret := gang{}
	for r, line := range lines {
		for x, c := range line {
			switch c {
			case '#':
				ret = append(ret, elf{pos{x, r}})
			}
		}
	}
	return ret
}

func (g gang) MoveN(n int) gang {
	moving := g.Copy()

	for elapsed := 0; elapsed < n; elapsed++ {
		current := make(map[pos]*elf)
		proposals := make(map[pos]*elf)

		for i := range moving {
			current[moving[i].loc] = &moving[i]
		}

		for j := range moving {
			e := &moving[j]
			p := e.GetProposal(current, elapsed)
			if p != nil {
				e.TryAddToMap(*p, proposals)
			}
		}

		moves := 0
		// fmt.Println("BEfore", moving)
		for k, v := range proposals {
			if v != POISON && v.loc != k {
				v.loc = k
				moves++
			}
		}
		// if moves == 0
		{
			fmt.Println("At step", elapsed+1, "moves", moves)
		}

		// fmt.Println("After ", moving)
		// fmt.Println("After Step", elapsed)
		// moving.PrintBoard()
	}

	return moving
}

func (e *elf) GetProposal(start map[pos]*elf, elapsed int) *pos {
	for i := 0; i < 4; i++ {
		if start[pos{e.loc.x - 1, e.loc.y - 1}] == nil &&
			start[pos{e.loc.x, e.loc.y - 1}] == nil &&
			start[pos{e.loc.x + 1, e.loc.y - 1}] == nil &&
			start[pos{e.loc.x - 1, e.loc.y + 1}] == nil &&
			start[pos{e.loc.x, e.loc.y + 1}] == nil &&
			start[pos{e.loc.x + 1, e.loc.y + 1}] == nil &&

			start[pos{e.loc.x - 1, e.loc.y}] == nil &&
			start[pos{e.loc.x + 1, e.loc.y}] == nil {
			return &e.loc
		}
		stepi := (i + elapsed) % 4
		switch stepi {
		case N:
			if start[pos{e.loc.x - 1, e.loc.y - 1}] == nil &&
				start[pos{e.loc.x, e.loc.y - 1}] == nil &&
				start[pos{e.loc.x + 1, e.loc.y - 1}] == nil {
				return &pos{e.loc.x, e.loc.y - 1}
			}
		case S:
			if start[pos{e.loc.x - 1, e.loc.y + 1}] == nil &&
				start[pos{e.loc.x, e.loc.y + 1}] == nil &&
				start[pos{e.loc.x + 1, e.loc.y + 1}] == nil {
				return &pos{e.loc.x, e.loc.y + 1}
			}
		case W:
			if start[pos{e.loc.x - 1, e.loc.y - 1}] == nil &&
				start[pos{e.loc.x - 1, e.loc.y}] == nil &&
				start[pos{e.loc.x - 1, e.loc.y + 1}] == nil {
				return &pos{e.loc.x - 1, e.loc.y}
			}
		case E:
			if start[pos{e.loc.x + 1, e.loc.y - 1}] == nil &&
				start[pos{e.loc.x + 1, e.loc.y}] == nil &&
				start[pos{e.loc.x + 1, e.loc.y + 1}] == nil {
				return &pos{e.loc.x + 1, e.loc.y}
			}
		default:
			fmt.Println("HOW?")
		}
	}
	// panic("what if blocked on all sides?")
	return nil
}

func (e *elf) TryAddToMap(newLoc pos, m map[pos]*elf) {
	if _, ok := m[newLoc]; ok {
		//Found existing
		m[newLoc] = POISON
	} else {
		m[newLoc] = e
	}
}

func (g gang) Copy() gang {
	ret := make(gang, len(g))
	copy(ret, g)
	return ret
}

func (g gang) CountSparseEmpties() int {
	NW, SE := g.GetBoundingBox()
	xLen := SE.x - NW.x + 1
	yLen := SE.y - NW.y + 1

	return xLen*yLen - len(g)
}

func (g gang) GetBoundingBox() (NW, SE pos) {
	NW = pos{math.MaxInt, math.MaxInt}
	SE = pos{math.MinInt, math.MinInt}
	for _, e := range g {
		if e.loc.x < NW.x {
			NW.x = e.loc.x
		}
		if e.loc.x > SE.x {
			SE.x = e.loc.x
		}
		if e.loc.y < NW.y {
			NW.y = e.loc.y
		}
		if e.loc.y > SE.y {
			SE.y = e.loc.y
		}
	}
	return NW, SE
}

func (g gang) PrintBoard() {
	NW, SE := g.GetBoundingBox()

	for y := NW.y; y <= SE.y; y++ {
		for x := NW.x; x <= SE.x; x++ {
			if g.HasElfAt(pos{x, y}) {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println()
	}
}

func (g gang) HasElfAt(loc pos) bool {
	for i := range g {
		if g[i].loc == loc {
			return true
		}
	}
	return false
}
