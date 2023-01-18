package main

import (
	_ "embed"
	"fmt"
)

var (
	//go:embed input.txt
	INPUT string

	//go:embed input2.txt
	INPUT2 string
)

type pos struct {
	x int
	y int
}

type rock []pos

type jets struct {
	seq string
	pos int
}

func main() {

	j := parseJetSequence(INPUT)
	j2 := parseJetSequence(INPUT2)

	fmt.Println("Part1 EXAMPLE", j2.DropRocks(2022))
	fmt.Println("Part1", j.DropRocks(2022))
}

func parseJetSequence(input string) jets {
	return jets{seq: input, pos: -1}
}

// DropRocks n times and return the highest resting ground/rock position
func (j jets) DropRocks(n int) int {
	floor := make(map[pos]bool)
	floorHeight := 0

	for i := 1; i <= n; i++ {
		r := NewRock(i, floorHeight+4)

		for {
			if j.NextIsLeft() {
				_ = r.ShiftXY(-1, 0, floor)
				// fmt.Println("Left", , r)
			} else {
				_ = r.ShiftXY(1, 0, floor)
				// fmt.Println("Right", , r)
			}

			if ok := r.ShiftXY(0, -1, floor); !ok {
				break
			}
		}
		floorHeight = max(floorHeight, r.MaxY())
		// fmt.Println("Can't fall further", r, floorHeight)
		for _, p := range r {
			floor[p] = true
		}

		// printTower(floor, floorHeight)
	}
	return floorHeight
}

func (j *jets) NextIsLeft() bool {
	j.pos++
	if j.pos >= len(j.seq) {
		j.pos = 0
	}
	return j.seq[j.pos] == '<'
}

func (r *rock) ShiftXY(xDelta, yDelta int, blocked map[pos]bool) bool {
	for i := range *r {
		if isBlocked(pos{(*r)[i].x + xDelta, (*r)[i].y + yDelta}, blocked) {
			return false
		}
	}
	for i := range *r {
		(*r)[i].x += xDelta
		(*r)[i].y += yDelta
	}
	return true
}

// Legal y > 0; x > 0 && x < 8
func isBlocked(p pos, b map[pos]bool) bool {
	if p.x < 1 || p.x > 7 {
		return true
	}
	if p.y < 1 {
		return true
	}
	if b[p] {
		return true
	}
	return false
}

func (r rock) MaxY() int {
	ret := 0
	for _, p := range r {
		if p.y > ret {
			ret = p.y
		}
	}
	return ret
}

// NewRock assumes a 1-based param,
func NewRock(i int, y0 int) rock {
	var ret rock
	if i == 0 {
		panic("we accept int with min 1")
	}
	switch (i - 1) % 5 {
	case 0:
		ret = rock{
			{3, y0},
			{4, y0},
			{5, y0},
			{6, y0},
		}
	case 1:
		ret = []pos{
			{4, y0},
			{3, y0 + 1},
			{4, y0 + 1},
			{5, y0 + 1},
			{4, y0 + 2},
		}
	case 2:
		ret = rock{
			{3, y0},
			{4, y0},
			{5, y0},
			{5, y0 + 1},
			{5, y0 + 2},
		}
	case 3:
		ret = rock{
			{3, y0},
			{3, y0 + 1},
			{3, y0 + 2},
			{3, y0 + 3},
		}
	case 4:
		ret = rock{
			{3, y0},
			{4, y0},
			{3, y0 + 1},
			{4, y0 + 1},
		}
	}
	return ret
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func printTower(blocked map[pos]bool, maxHeight int) {
	for y := maxHeight; y > 0; y-- {
		fmt.Print("|")
		for x := 1; x <= 7; x++ {
			if blocked[pos{x, y}] {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("|")
	}
}
