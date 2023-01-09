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

type pos struct {
	x int
	y int
}

type player struct {
	loc pos
	dir heading
}

type heading int8

type rowDescr struct {
	minX int
	maxX int
}

type colDescr struct {
	minY int
	maxY int
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
	dir   heading
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

		for i := 0; i < sc.moves; i++ {
			cur.dir = sc.dir
			cur.loc.NextStep(sc.dir, g)

		}

	}

	return -1, cur
}

func (c command) GetSubcommands() []subcommand {
	ret := make([]subcommand, 0)

	//First Move has no rotation
	dir := RIGHT
	sc := subcommand{moves: 0, dir: dir}
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
			ret = append(ret, parseSubcommand(str, dir))
			str = "" //We're Done
		} else {
			nextSCStart++
			sc := parseSubcommand(str[:nextSCStart], dir)
			dir = sc.dir
			ret = append(ret, sc)
			str = str[nextSCStart:]
		}

	}
	return ret
}

func parseSubcommand(str string, prev heading) subcommand {
	ret := subcommand{}
	switch str[0] {
	case 'L':
		ret.dir = prev.Rotate(false)
	case 'R':
		ret.dir = prev.Rotate(true)
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

func (h heading) Rotate(clockwise bool) heading {
	if clockwise {
		n := (h + 1) % 4
		return n
	}

	n := h - 1
	if n < 0 {
		n = 3
	}
	return n
}

func (p *pos) NextStep(dir heading, g *grid) {
	switch dir {
	case UP:
		nextY := p.y - 1
		if nextY < g.cols[p.x].minY {
			nextY = g.cols[p.x].maxY
		}

		if g.spaces[nextY][p.x] {
			p.y = nextY
		}

	case DOWN:
		nextY := p.y + 1
		if nextY > g.cols[p.x].maxY {
			nextY = g.cols[p.x].minY
		}

		if g.spaces[nextY][p.x] {
			p.y = nextY
		}
	case LEFT:
		nextX := p.x - 1
		if nextX < g.rows[p.y].minX {
			nextX = g.rows[p.y].maxX
		}

		if g.spaces[p.y][nextX] {
			p.x = nextX
		}
	case RIGHT:
		nextX := p.x + 1
		if nextX > g.rows[p.y].maxX {
			nextX = g.rows[p.y].minX
		}

		if g.spaces[p.y][nextX] {
			p.x = nextX
		}
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
