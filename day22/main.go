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
	// fmt.Println(lines)

	{
		g, cmd := parseInput(lines, false)
		start := player{loc: pos{g.rows[0].minX, 0}, dir: RIGHT}
		_, end := cmd.RunCommand(&g, start)

		row, col := end.loc.y+1, end.loc.x+1
		fmt.Println("Finished at", col, row, "facing", end.dir)
		fmt.Println("Part1 Score", 1000*row+4*col+end.dir.Score())
		//Remember to print row, col with "+1" (start at 1)
	}

	{
		g2, cmd := parseInput(lines, true)
		start := player{loc: pos{g2.rows[0].minX, 0}, dir: RIGHT}
		_, end := cmd.RunCommand(&g2, start)

		row, col := end.loc.y+1, end.loc.x+1
		fmt.Println("Finished at", col, row, "facing", end.dir)
		fmt.Println("Part1 Score", 1000*row+4*col+end.dir.Score())

		//Remember to print row, col with "+1" (start at 1)
	}
}

func parseInput(lines []string, part2 bool) (grid, command) {
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

	totalSpaces := 0
	g.rows = make([]rowDescr, numRows)
	for y := 0; y < numRows; y++ {
		rd := rowDescr{minX: math.MaxInt, maxX: 0}
		for x := 0; x < numCols; x++ {
			if len(board[y]) > x && board[y][x] != ' ' {
				rd.minX = min(rd.minX, x)
				rd.maxX = max(rd.maxX, x)
				totalSpaces++
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

	if part2 {
		cubeHeight := int(math.Sqrt(float64(totalSpaces) / 6))
		fmt.Println("total spaces", totalSpaces, "Cube Height", cubeHeight)

		if cubeHeight != 50 {
			panic("I only hard coded folding for the 50")
		}
		for y := 0; y < 50; y++ {
			// Tile 1 MinX Left goes to Tile 4 Right
			rd := &g.rows[y]
			foldedRow := &g.rows[149-y]

			rd.minXNext = player{loc: pos{foldedRow.minX, 149 - y}, dir: LEFT.RotateN(COUNTERCLOCKWISE, 2)}
			//do the reverse, Left MinX from Tile4 goes Right on Tile 1
			foldedRow.minXNext = player{loc: pos{rd.minX, y}, dir: LEFT.RotateN(CLOCKWISE, 2)}

			//Tile 2 MaxR Right goes to Tile 5 left
			rd.maxXNext = player{loc: pos{foldedRow.maxX, 149 - y}, dir: RIGHT.RotateN(COUNTERCLOCKWISE, 2)}
			//Do the reverse, Right from Tile 5 goes Left on Tile 2
			foldedRow.maxXNext = player{loc: pos{rd.maxX, y}, dir: RIGHT.RotateN(CLOCKWISE, 2)}
		}
		for y := 50; y < 100; y++ {

			//Tile 3 MaxX Right goes to Tile 2 Up
			rd := &g.rows[y]
			maxCol := &g.cols[y+50]
			rd.maxXNext = player{loc: pos{y + 50, maxCol.maxY}, dir: RIGHT.RotateN(COUNTERCLOCKWISE, 1)}

			//do the reverse, down from tile 2 goes left tile 3
			maxCol.maxYNext = player{loc: pos{rd.maxX, y}, dir: DOWN.RotateN(CLOCKWISE, 1)}

			//Left on Tile 3 goes Down on Tile 4
			minCol := &g.cols[y-50]
			rd.minXNext = player{loc: pos{y - 50, minCol.minY}, dir: LEFT.RotateN(COUNTERCLOCKWISE, 1)}

			//Do the reverse, going up on tile 4 goes right on tile 3
			minCol.minYNext = player{loc: pos{rd.minX, y}, dir: UP.RotateN(CLOCKWISE, 1)}
		}
		// Tiless 4 and 5 are already MinX/MaxX mapped

		for y := 150; y < 199; y++ {
			// Tile 6 MaxX Right goes to Tile 5 Up
			rd := &g.rows[y]
			maxCol := &g.cols[y-100]
			rd.maxXNext = player{loc: pos{y - 100, maxCol.maxY}, dir: RIGHT.RotateN(COUNTERCLOCKWISE, 1)}

			// do the reverse, down from tile 5 goes left on tile 6
			maxCol.maxYNext = player{loc: pos{rd.maxX, y}, dir: DOWN.RotateN(CLOCKWISE, 1)}

			//Left on Tile 6 goes Down on Tile 1
			minCol := &g.cols[y-100]
			rd.minXNext = player{loc: pos{y - 100, minCol.minY}, dir: LEFT.RotateN(COUNTERCLOCKWISE, 1)}

			//Do the reverse, going up on tile 4 goes right on tile 3
			minCol.minYNext = player{loc: pos{rd.minX, y}, dir: UP.RotateN(CLOCKWISE, 1)}
		}

		for x := 50; x < 99; x++ {
			//Tile1 UP MinY turns into Tile 6 Right
			cd := &g.cols[x]

			maxRow := &g.rows[x+100]
			cd.minYNext = player{loc: pos{maxRow.minX, x + 100}, dir: UP.RotateN(COUNTERCLOCKWISE, 3)}

			//do the reverse, Tile 6 left MinX turns into Tile 1 Down
			maxRow.maxXNext = player{loc: pos{x, cd.minY}, dir: LEFT.RotateN(CLOCKWISE, 3)}

		}

		for x := 0; x < 50; x++ {
			//Tile 6 down MaxY turns into Tile 2 down
			cd := &g.cols[x]

			colMap := &g.cols[x+100]
			cd.maxYNext = player{loc: pos{x + 100, colMap.minY}, dir: DOWN}
			//Reverse - Tile2 Up turns into Tile 6 up
			cd.minYNext = player{loc: pos{x, colMap.maxY}, dir: UP}

		}

	} else {
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

func (h heading) RotateN(dir rotation, n int) heading {
	ret := h
	for i := 0; i < n; i++ {
		ret = ret.Rotate(dir)
	}
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
