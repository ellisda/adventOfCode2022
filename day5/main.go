package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	stacks, moves := parseInput(f)
	// fmt.Println("stacks", stacks)
	// printStacks(stacks)

	// ProcessMoves9000(stacks, moves)

	ProcessMoves9001(stacks, moves)

	// bernard()
}

type stack struct {
	items []rune
}

type move struct {
	fromCol int
	toCol   int
	num     int
}

func ProcessMoves9000(in []*stack, moves []move) {
	stacks := deepCopy(in)

	for _, m := range moves {
		for i := 0; i < m.num; i++ {
			stacks[m.toCol-1].Push(stacks[m.fromCol-1].Pop())
		}
		// fmt.Printf("After move %d from %d to %d \n", m.num, m.fromCol, m.toCol)
		// printStacks(stacks)

	}
	// fmt.Println("After")
	// printStacks(stacks)

	fmt.Println("9000 Tops: ", GetTops(stacks))
}

func ProcessMoves9001(in []*stack, moves []move) {
	stacks := deepCopy(in)

	for _, m := range moves {
		stacks[m.toCol-1].PushN(stacks[m.fromCol-1].PopN(m.num))
		// fmt.Printf("After move %d from %d to %d \n", m.num, m.fromCol, m.toCol)
		// printStacks(stacks)

	}
	fmt.Println("9001 Tops: ", GetTops(stacks))

	fmt.Println("After")
	printStacks(stacks)
}

func ProcessMoves9001Bad(in []*stack, moves []move) {
	stacks := deepCopy(in)

	for _, m := range moves {
		popped := make([]rune, m.num)
		for i := 0; i < m.num; i++ {
			popped[i] = stacks[m.fromCol-1].Pop()
		}

		for i := len(popped) - 1; i >= 0; i-- {
			stacks[m.toCol-1].Push(popped[i])
		}
		// fmt.Printf("After move %d from %d to %d \n", m.num, m.fromCol, m.toCol)
		// printStacks(stacks)

	}
	// fmt.Println("After")
	// printStacks(stacks)
	// fmt.Println("9001 Tops: ", GetTops(stacks))
}

func GetTops(stacks []*stack) string {
	ret := ""
	for _, s := range stacks {
		ret = ret + fmt.Sprintf("%c", s.Peek())
	}
	return ret
}

func (s stack) Len() int {
	return len(s.items)
}

func (s *stack) Push(r rune) {
	s.items = append(s.items, r)
}

// Push N where last input will be top of stack
func (s *stack) PushN(r []rune) {
	s.items = append(s.items, r...)
}

func (s *stack) Pop() rune {
	r := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]

	return r
}

// Top of stack will be last item in returned slice
func (s *stack) PopN(n int) []rune {
	r := s.items[len(s.items)-n : len(s.items)]
	s.items = s.items[:len(s.items)-n]

	return r
}

func (s *stack) Peek() rune {
	return s.items[len(s.items)-1]
}

func parseInput(r io.Reader) ([]*stack, []move) {
	//4--char colums, where we want values from index 1, 5, 9, etc.

	lines := []string{}
	pictureAboveIdx := -1

	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		if pictureAboveIdx == -1 && len(line) == 0 {
			pictureAboveIdx = len(lines)
		}
		lines = append(lines, line)
	}

	stacks := make([]*stack, len(lines[0])/4+1)
	for i := range stacks {
		stacks[i] = &stack{}
	}

	for i := pictureAboveIdx - 2; i >= 0; i-- {
		for col := 1; col < len(lines[i]); col += 4 {
			if lines[i][col] != ' ' {
				stack := stacks[col/4]
				stack.Push(rune(lines[i][col]))
				// stacks[col/4] = stack
			}
		}
	}

	moves := []move{}
	for i := pictureAboveIdx + 1; i < len(lines); i++ {
		m := move{}
		n, err := fmt.Sscanf(lines[i], "move %d from %d to %d", &m.num, &m.fromCol, &m.toCol)
		if n != 3 || err != nil {
			log.Fatal("Expected to parse 3 items, got ", n, "error:", err)
		}
		moves = append(moves, m)
	}
	return stacks, moves

}

func deepCopy(in []*stack) []*stack {
	ret := make([]*stack, len(in))

	n := copy(ret, in)
	if n < len(in) {
		panic("failed to copy input")
	}
	return ret
}

func printStacks(stacks []*stack) {
	//Take a deep copy, so I can desctrutively Pop() items to print
	// stacks := deepCopy(in)
	itemsLeft := 0
	for _, s := range stacks {
		if itemsLeft < s.Len() {
			itemsLeft = s.Len()
		}
	}

	fmt.Println("Stacks:")
	for ; itemsLeft > 0; itemsLeft-- {
		for si, s := range stacks {
			if s.Len() == itemsLeft {

				r := stacks[si].Pop()
				// r := stacks[si].items[itemsLeft-1]
				fmt.Printf("[%c] ", r)
			} else {
				fmt.Print("    ")
			}
		}
		fmt.Println()
	}
}
