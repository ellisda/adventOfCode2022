package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
)

func ProcessMoves9000_Old(in []stack, moves []move) {
	stacks := deepCopy_old(in)

	for _, m := range moves {
		for i := 0; i < m.num; i++ {
			stacks[m.toCol-1].Push(stacks[m.fromCol-1].Pop())
		}
		// fmt.Printf("After move %d from %d to %d \n", m.num, m.fromCol, m.toCol)
		// printStacks(stacks)

	}
	fmt.Println("After")
	printStacks_old(stacks)

	// fmt.Println("9000 Tops: ", GetTops(stacks))
}

func ProcessMoves9001_Old(in []stack, moves []move) {
	stacks := deepCopy_old(in)

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
	fmt.Println("After")
	printStacks_old(stacks)
	// fmt.Println("9001 Tops: ", GetTops(stacks))
}

func parseInput_old(r io.Reader) ([]stack, []move) {
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

	stacks := make([]stack, len(lines[0])/4+1)
	for i := pictureAboveIdx - 2; i >= 0; i-- {
		for col := 1; col < len(lines[i]); col += 4 {
			if lines[i][col] != ' ' {
				stack := &stacks[col/4]
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

func deepCopy_old(in []stack) []stack {
	ret := make([]stack, len(in))

	for i := 0; i < len(in); i++ {
		s_in := in[i]

		s := stack{items: make([]rune, s_in.Len())}
		// for j := 0; j < s_in.Len(); j++ {
		n := copy(s.items, s_in.items)
		if n < in[i].Len() {
			panic("failed to copy input")
			// }
		}
		ret[i] = s
	}
	return ret
}

func printStacks_old(in []stack) {
	//Take a deep copy, so I can desctrutively Pop() items to print
	stacks := deepCopy_old(in)
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
				fmt.Printf("[%c] ", r)
			} else {
				fmt.Print("    ")
			}
		}
		fmt.Println()
	}
}
