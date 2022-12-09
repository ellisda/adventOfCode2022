package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	Up    direction = 'U'
	Down  direction = 'D'
	Left  direction = 'L'
	Right direction = 'R'
)

type position struct {
	x int
	y int
}

func (p *position) Move(d direction) {
	switch d {
	case Up:
		p.y++
	case Down:
		p.y--
	case Left:
		p.x--
	case Right:
		p.x++
	}
}

type move struct {
	dir      direction
	distance int
}

type direction rune

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	moves := parseInput(f)
	tailPositions := processMoves(moves, 2)

	fmt.Println("part 1", len(tailPositions))

	fmt.Println("part 2", len(processMoves(moves, 10)))

}

func processMoves(moves []move, numKnots int) (tailPositions map[position]int) {
	tailPositions = make(map[position]int)
	knots := make([]position, numKnots)
	for _, m := range moves {
		fmt.Print("Beginning ", m)
		for i := 0; i < m.distance; i++ {
			knots[0].Move(m.dir)

			for i := 1; i < len(knots); i++ {
				tailCatchup(knots[i-1], &knots[i])
			}
			tailPositions[knots[len(knots)-1]]++
			fmt.Print(string(m.dir), knots[len(knots)-1], "total", len(tailPositions))
		}
		fmt.Println()
	}

	return
}

func tailCatchup(head position, tail *position) {
	xDelta, yDelta := head.x-tail.x, head.y-tail.y
	if abs(xDelta) <= 1 && abs(yDelta) <= 1 {
		return
	}

	if yDelta != 0 {
		tail.y += (yDelta / abs(yDelta)) // unit Move that preserves sign
	}
	if xDelta != 0 {
		tail.x += (xDelta / abs(xDelta)) // unit Move that preserves sign
	}
}

func abs(n int) int {
	if n >= 0 {
		return n
	}
	return -1 * n
}

func parseInput(f io.ReadSeekCloser) []move {

	ret := []move{}
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()

		m := move{}
		fmt.Sscanf(line, "%c %d", &m.dir, &m.distance)

		ret = append(ret, m)
	}
	return ret
}
