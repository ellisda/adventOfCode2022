package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	UP heading = iota
	RIGHT
	DOWN
	LEFT
)

const (
	CLOCKWISE rotation = iota
	COUNTERCLOCKWISE
	NOROTATION
)

type pos struct {
	x int
	y int
}

type player struct {
	loc pos
	dir heading
}

type heading int8

// rowDescr describes the min/max range for valid x values in the row
// and what next position and heading to go to if you walk out of range
type rowDescr struct {
	minX     int
	minXNext player

	maxX     int
	maxXNext player
}

type colDescr struct {
	minY     int
	minYNext player

	maxY     int
	maxYNext player
}

// a freespace is true if a player can stand there
type freespace bool

type grid struct {
	rows   []rowDescr
	cols   []colDescr
	spaces [][]freespace
}

type command string

type subcommand struct {
	moves int
	dir   rotation
}
type rotation int8

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
	fmt.Println(lines)

	g, cmd := parseInput(lines)
	fmt.Println(cmd.GetSubcommands())

	start := player{loc: pos{g.rows[0].minX, 0}, dir: RIGHT}
	_, end := cmd.RunCommand(&g, start)

	row, col := end.loc.y+1, end.loc.x+1
	fmt.Println("Finished at", col, row, "facing", end.dir)
	fmt.Println("Part1 Score", 1000*row+4*col+end.dir.Score())

	//Remember to print row, col with "+1" (start at 1)
}

func parseInput(lines []string) (grid, command) {
	board := lines[:len(lines)-2]

	numRows := len(board)
	numCols := len(board[0])

	g := grid{}
	g.spaces = make([][]freespace, numRows)

	for _, line := range board {
		if len(line) > numCols {
			numCols = len(line)
		}
	}

	for r, line := range board {
		row := make([]freespace, numCols)
		g.spaces[r] = row
		for x, c := range line {
			if c == '.' {
				row[x] = true
			}
		}
	}

	g.rows = make([]rowDescr, numRows)
	for y := 0; y < numRows; y++ {
		rd := rowDescr{minX: math.MaxInt, maxX: 0}
		for x := 0; x < numCols; x++ {
			if len(board[y]) > x && board[y][x] != ' ' {
				rd.minX = min(rd.minX, x)
				rd.maxX = max(rd.maxX, x)
			}
		}
		g.rows[y] = rd
	}

	g.cols = make([]colDescr, numCols)
	for x := 0; x < numCols; x++ {
		cd := colDescr{minY: math.MaxInt, maxY: 0}
		for y := 0; y < numRows; y++ {
			if len(board[y]) > x && board[y][x] != ' ' {
				cd.minY = min(cd.minY, y)
				cd.maxY = max(cd.maxY, y)
			}
		}
		g.cols[x] = cd
	}

	//Part1
	for y := range g.rows {
		rd := &g.rows[y]
		rd.maxXNext = player{loc: pos{rd.minX, y}, dir: RIGHT} //same dir, but start over at MinX
		rd.minXNext = player{loc: pos{rd.maxX, y}, dir: LEFT}
	}
	for x := range g.cols {
		cd := &g.cols[x]
		cd.maxYNext = player{loc: pos{x, cd.minY}, dir: DOWN} //same dir, but start over at MinX
		cd.minYNext = player{loc: pos{x, cd.maxY}, dir: UP}
	}

	cmd := command(lines[len(lines)-1])
	return g, cmd
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (c command) RunCommand(g *grid, start player) (spaces int, final player) {
	cur := start
	for _, sc := range c.GetSubcommands() {
		cur.dir = cur.dir.Rotate(sc.dir)

		for i := 0; i < sc.moves; i++ {
			cur.NextStep(g)

		}

	}

	return -1, cur
}

func (c command) GetSubcommands() []subcommand {
	ret := make([]subcommand, 0)

	//First Move has no rotation
	sc := subcommand{moves: 0, dir: NOROTATION}
	i := strings.IndexAny(string(c), "LR")
	v, err := strconv.Atoi(string(c)[:i])
	if err != nil {
		log.Fatal("cmd parse fail")
	}
	sc.moves = v
	ret = append(ret, sc)

	str := string(c)[i:]
	for {
		if len(str) <= 1 {
			break
		}

		if nextSCStart := strings.IndexAny(str[1:], "LR"); nextSCStart == -1 {
			ret = append(ret, parseSubcommand(str))
			str = "" //We're Done
		} else {
			nextSCStart++
			sc := parseSubcommand(str[:nextSCStart])
			ret = append(ret, sc)
			str = str[nextSCStart:]
		}

	}
	return ret
}

func parseSubcommand(str string) subcommand {
	ret := subcommand{}
	switch str[0] {
	case 'L':
		ret.dir = COUNTERCLOCKWISE
	case 'R':
		ret.dir = CLOCKWISE
	default:
		panic("boo")
	}

	v, err := strconv.Atoi(str[1:])
	if err != nil {
		log.Fatal("cmd parse fail")
	}
	ret.moves = v
	return ret
}

func (h heading) Rotate(dir rotation) heading {
	switch dir {
	case NOROTATION:
		return h
	case CLOCKWISE:
		n := (h + 1) % 4
		return n
	case COUNTERCLOCKWISE:
		n := h - 1
		if n < 0 {
			n = 3
		}
		return n
	}
	panic("eek")
}

func (p *player) NextStep(g *grid) {
	next := *p
	switch next.dir {
	case UP:
		next.loc.y--
		if next.loc.y < g.cols[next.loc.x].minY {
			next = g.cols[p.loc.x].minYNext
		}
	case DOWN:
		next.loc.y++
		if next.loc.y > g.cols[next.loc.x].maxY {
			next = g.cols[next.loc.x].maxYNext
		}
	case LEFT:
		next.loc.x--
		if next.loc.x < g.rows[next.loc.y].minX {
			next = g.rows[next.loc.y].minXNext
		}
	case RIGHT:
		next.loc.x++
		if next.loc.x > g.rows[next.loc.y].maxX {
			next = g.rows[next.loc.y].maxXNext
		}
	}

	//If next position isn't blocked, go there, and update heading if we've walked onto another face
	if g.spaces[next.loc.y][next.loc.x] {
		*p = next
	}
}

func (h heading) Score() int {
	switch h {
	case UP:
		return 3
	case RIGHT:
		return 0
	case DOWN:
		return 1
	case LEFT:
		return 2
	}
	panic("boo")
}
