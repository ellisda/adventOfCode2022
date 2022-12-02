package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filename := "input.txt"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	games := parseInput(f)
	fmt.Println("found: ", games)

	total := 0
	for _, g := range games {
		fmt.Println(g, g.Score())
		total += g.Score()
	}
	fmt.Println("Total Score:", total)

	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	hints := parseInput2(f)

	total = 0
	for _, h := range hints {
		// fmt.Println("hint:", h, h.Game(), h.Game().Score())
		total += h.Game().Score()
	}
	fmt.Println("Total score from hint file: ", total)
}

const (
	RockThem     theirMove = "A"
	PaperThem    theirMove = "B"
	ScissorsThem theirMove = "C"

	RockMe     myMove = "X"
	PaperMe    myMove = "Y"
	ScissorsMe myMove = "Z"

	Lose gameResult = "X"
	Draw gameResult = "Y"
	Win  gameResult = "Z"
)

type theirMove string
type myMove string
type gameResult string

type game struct {
	them theirMove
	mine myMove
}

type hint struct {
	them   theirMove
	result gameResult
}

// Game returns the game struct with moves to produce the hint's result
func (h hint) Game() game {
	g := game{
		them: h.them,
	}
	switch h.result {
	case Win:
		switch h.them {
		case RockThem:
			g.mine = PaperMe
		case PaperThem:
			g.mine = ScissorsMe
		case ScissorsThem:
			g.mine = RockMe
		}
	case Lose:
		switch h.them {
		case RockThem:
			g.mine = ScissorsMe
		case PaperThem:
			g.mine = RockMe
		case ScissorsThem:
			g.mine = PaperMe
		}
	case Draw:
		switch h.them {
		case RockThem:
			g.mine = RockMe
		case PaperThem:
			g.mine = PaperMe
		case ScissorsThem:
			g.mine = ScissorsMe
		}
	}
	return g
}

// Score returns your myMove score plus 3 for a draw, 6 for a win, and 0 for a loss
func (g game) Score() int {
	winLossPoints := 0
	switch {
	case (g.mine == RockMe && g.them == ScissorsThem) ||
		(g.mine == PaperMe && g.them == RockThem) ||
		(g.mine == ScissorsMe && g.them == PaperThem):
		winLossPoints = 6 //Win
	case (g.mine == RockMe && g.them == RockThem) ||
		(g.mine == PaperMe && g.them == PaperThem) ||
		(g.mine == ScissorsMe && g.them == ScissorsThem):
		winLossPoints = 3 //Draw
	default:
		winLossPoints = 0 //Loss
	}

	return g.mine.Score() + winLossPoints
}

// Score returns 1 when you throw Rock (X), 2 for Paper (Y), and 3 for Scissors (Z)
func (m myMove) Score() int {
	return int(m[0]) - int('W')
}

func (h game) String() string {
	return fmt.Sprintf("{them: %s, me: %s}", h.them, h.mine)
}

func (m myMove) String() string {
	switch m {
	case RockMe:
		return "Rock"
	case PaperMe:
		return "Paper"
	case ScissorsMe:
		return "Scissors"
	default:
		panic("unknown move type " + m)
	}
}
func (t theirMove) String() string {
	switch t {
	case RockThem:
		return "Rock"
	case PaperThem:
		return "Paper"
	case ScissorsThem:
		return "Scissors"
	default:
		panic("unknown move type " + t)
	}
}

func parseInput(r io.Reader) []game {
	ret := []game{}
	for {
		g := game{}
		n, err := fmt.Fscanln(r, &g.them, &g.mine)

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}

		ret = append(ret, g)
	}

	return ret
}

func parseInput2(r io.Reader) []hint {
	ret := []hint{}
	for {
		h := hint{}
		n, err := fmt.Fscanln(r, &h.them, &h.result)

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}

		ret = append(ret, h)
	}

	return ret
}
