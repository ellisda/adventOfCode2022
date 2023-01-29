package main

import (
	"bufio"
	"container/heap"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

type blizzard int8

// type player int8

const (
	Up blizzard = iota
	Down
	Left
	Right
	Multiple
	Empty
	// player player = -1
)

type blizzardStart struct {
	start position
	dir   blizzard

	X_len *int //all blizzards are on same board, and share max coordinates
	Y_len *int
}

type position struct {
	x int
	y int
}

// a simulation is a collection of known blizard starting locations and directions
type simulation []blizzardStart

// // a stormgrid is list of storm positions for a given point in time
// type stormgrid [][]bool

type move struct {
	from    position
	to      position
	stepsTo int
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)

		t0 := time.Now()

		go func() {
			for {
				if time.Since(t0) > 10*time.Second {
					pprof.StopCPUProfile()
					panic("timeout")
				}
				time.Sleep(50 * time.Millisecond)
			}
		}()
	}

	in := "input.txt"
	file, _ := os.Open(in)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	sim, start, end := ParseStorms(lines)
	fmt.Println("start", start, "end", end)
	fmt.Println("sim", sim)

	fmt.Println("Part1", sim.WalkBFS_Priority(start, end, 500))
}

func ParseStorms(lines []string) (ret simulation, start, end position) {
	xLen, yLen := len(lines[0])-2, len(lines)-2

	// Coordinates exclude the border, so interior space in uppoer left is 0,0 and both start and end are outside normal bounds
	for r, line := range lines {
		for x, c := range line {
			switch c {
			case '#':
			case '>':
				ret = append(ret, blizzardStart{position{x - 1, r - 1}, Right, &xLen, &yLen})
			case '<':
				ret = append(ret, blizzardStart{position{x - 1, r - 1}, Left, &xLen, &yLen})
			case 'v':
				ret = append(ret, blizzardStart{position{x - 1, r - 1}, Down, &xLen, &yLen})
			case '^':
				ret = append(ret, blizzardStart{position{x - 1, r - 1}, Up, &xLen, &yLen})
			case '.':
			}
		}
	}
	// Start position is at -1y, and end position is one y past bottom of normal dimensions
	start = position{x: strings.IndexRune(lines[0], '.') - 1, y: -1}
	end = position{x: strings.IndexRune(lines[len(lines)-1], '.') - 1, y: len(lines) - 2}
	return ret, start, end
}

func (s *simulation) WalkBFS_Priority(start, end position, maxSteps int) int {
	q := &IntHeap{dest: end}

	heap.Push(q, move{from: start, to: position{start.x, 0}, stepsTo: 1})

	copy := s.Copy()
	breadth := 0
	for {
		m := heap.Pop(q).(move)
		if m.to == end {
			return m.stepsTo
		}
		if m.stepsTo > breadth {
			breadth = m.stepsTo
			log.Println("Evaluating Moves with NumSteps", breadth, "queue len", q.Len(), "starting with", m)
		}
		if m.stepsTo > maxSteps {
			log.Fatal("Exiting. Popped a move that has already taken maxSteps", m)
		}

		for _, n := range s.GetFreeSpaces(copy, m.to, m.stepsTo+1, end) {
			//Walk each
			// do we ever stop if we've been to this position 30+ times already?
			// A: no, we just priorize moves that get us closer to the target

			heap.Push(q, move{from: m.to, to: n, stepsTo: m.stepsTo + 1})
		}

	}

}

// // WalkOne position to next, in the goal of reaching the target. Returns number of steps taken
// // at time that it finally reaches destination, or -1 if destination was not reached
// func (s *simulation) WalkOne(target position, next position, stepsTaken int) int {
// 	if next == target {
// 		return stepsTaken
// 	}

// 	//Player should now look at the surrounding 9 positions and see which if any moves are possible
// 	for _, n := range s.GetFreeSpaces(next, stepsTaken) {
// 		//Walk each
// 		// do we ever stop if we've been to this position 30+ times already?
// 		score := s.WalkOne(target, n, stepsTaken+1)
// 		_ = score
// 	}

// 	return -1
// }

func (s *simulation) GetFreeSpaces(copy *simulation, player position, elapsed int, finalDest position) []position {
	// MEMORY INTENSIVE
	// DESIGN REVIEW - We don't really need to generate a whole grid, to look at 8 positions.
	// 				   instead, we just run each storm movement and check each to see if they are
	//                 present in the 8 positions
	today := s.RunBlizzardsFast(copy, elapsed)

	ret := make([]position, 0, 5)

	//DESIGN REVIEW - player is alowwed to stay put if no storm is coming into their current location
	for x := player.x - 1; x <= player.x+1; x++ {
		xPos := position{x, player.y}
		if xPos == finalDest || (x >= 0 && x < *(*s)[0].X_len &&
			!today.IsStorming(xPos)) {
			ret = append(ret, xPos)
		}
	}

	for y := player.y - 1; y <= player.y+1; y += 2 {
		yPos := position{player.x, y}
		if yPos == finalDest || (y >= 0 && y < *(*s)[0].Y_len &&
			!today.IsStorming(yPos)) {
			ret = append(ret, yPos)
		}
	}
	return ret
}

// Runs the initial blzzard simulation through to n number of steps and returns a grid
func (s *simulation) RunBlizzards(elapsed int) *simulation {
	ret := s.Copy()
	return s.RunBlizzardsFast(ret, elapsed)
}

func (s *simulation) RunBlizzardsFast(ret *simulation, elapsed int) *simulation {
	for i, b := range *s {
		x := &((*ret)[i])
		*x = b
		switch b.dir {
		case Down:
			x.start.y = (b.start.y + elapsed) % *b.Y_len
		case Up:
			x.start.y = (b.start.y - elapsed) % *b.Y_len
			if x.start.y < 0 {
				x.start.y += *b.Y_len
			}

		case Right:
			x.start.x = (b.start.x + elapsed) % *b.X_len
		case Left:
			x.start.x = (b.start.x - elapsed) % *b.X_len
			if x.start.x < 0 {
				x.start.x += *b.X_len
			}
		}

	}

	return ret
}

func (s *simulation) Copy() *simulation {
	ret := make(simulation, len(*s))
	copy(ret, *s)
	return &ret
}

func (s *simulation) IsStorming(p position) bool {
	for _, b := range *s {
		if b.start.x == p.x && b.start.y == p.y {
			return true
		}
	}
	return false
}
