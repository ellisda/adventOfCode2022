package main

import (
	_ "embed"
	"fmt"
	"math"
	"strconv"
	"strings"
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
	z int
}

type surface []pos

func main() {

	s2 := parseInput(strings.Split(INPUT2, "\n"))
	s := parseInput(strings.Split(INPUT, "\n"))

	fmt.Println("Part1 (EXAMPLE)", s2.CountFreeSides())
	fmt.Println("Part1", s.CountFreeSides())

	fmt.Println("Part2 (EXAMPLE)", s2.OutwardFacingSurfaces())
	fmt.Println("Part2", s.OutwardFacingSurfaces())
}

func parseInput(lines []string) surface {
	s := make(surface, len(lines))
	for i, line := range lines {
		xyz := strings.Split(line, ",")
		s[i] = pos{par(xyz[0]), par(xyz[1]), par(xyz[2])}
	}
	return s
}

func par(s string) int {
	ret, _ := strconv.Atoi(s)
	return ret
}

func (s surface) CountFreeSides() int {
	ret := 0
	for i := range s {
		for _, n := range s[i].Neighbors() {
			if s.IsEmpty(n) {
				ret++
			}
		}
	}
	return ret
}

func (p pos) Neighbors() []pos {
	return []pos{
		{p.x - 1, p.y, p.z},
		{p.x, p.y - 1, p.z},
		{p.x, p.y, p.z - 1},
		{p.x + 1, p.y, p.z},
		{p.x, p.y + 1, p.z},
		{p.x, p.y, p.z + 1},
	}
}

func (s surface) IsEmpty(p pos) bool {
	for i := range s {
		if s[i] == p {
			return false
		}
	}
	return true
}

func (s surface) OutwardFacingSurfaces() int {
	water := FillWithWater(s)
	ret := 0
	for i := range s {
		for _, n := range s[i].Neighbors() {
			if water[n] {
				ret++
			}
		}
	}
	return ret
}

func FillWithWater(s surface) map[pos]bool {
	p0, p1 := s.BoundingBox()
	ret := make(map[pos]bool)

	s.visit(p0, ret, p0, p1)
	return ret
}

func (s surface) BoundingBox() (min, max pos) {
	minAll := math.MaxInt
	maxAll := math.MinInt
	for i := range s {
		minAll = Min(Min(Min(minAll, s[i].x), s[i].y), s[i].z)
		maxAll = Max(Max(Max(maxAll, s[i].x), s[i].y), s[i].z)

	}
	minAll -= 1
	maxAll += 1
	return pos{minAll, minAll, minAll}, pos{maxAll, maxAll, maxAll}
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (s surface) visit(p pos, water map[pos]bool, min, max pos) {
	if _, ok := water[p]; ok {
		return
	}

	water[p] = true

	for _, neighbor := range p.Neighbors() {
		if inBounds(neighbor, min, max) && s.IsEmpty(neighbor) {
			s.visit(neighbor, water, min, max)
		}
	}
}

func inBounds(p, min, max pos) bool {
	out := p.x < min.x || p.x > max.x ||
		p.y < min.y || p.y > max.y ||
		p.z < min.z || p.z > max.z

	return !out
}
