package main

import (
	"fmt"
	"io"
	"log"
	"os"
)

type pair struct {
	p1Start, p1Stop int
	p2Start, p2Stop int
}

func main() {
	filename := "input.txt"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	pairs := parseInput(f)
	fmt.Println("found: ", pairs)

	total := 0
	for _, p := range pairs {
		if isContainedBy(p) {
			fmt.Println("Contained", p)
			total++
		}
	}
	fmt.Println("Total", total)

	problem2(pairs)
}

func problem2(pairs []pair) {
	total := 0
	for _, p := range pairs {
		if !isDisjoint(p) {
			fmt.Println("Overlapped", p)
			total++
		}
	}
	fmt.Println("Total", total)

}

func isContainedBy(p pair) bool {
	return (p.p1Start <= p.p2Start && p.p1Stop >= p.p2Stop) ||
		(p.p2Start <= p.p1Start && p.p2Stop >= p.p1Stop)
}

func hasOverlap(p pair) bool {
	return (p.p1Stop >= p.p2Start && p.p1Start <= p.p2Start) ||
		(p.p2Stop >= p.p1Start && p.p2Start <= p.p1Start)
}

func isDisjoint(p pair) bool {
	return p.p1Start > p.p2Stop || p.p1Stop < p.p2Start
}

func parseInput(r io.Reader) []pair {
	ret := []pair{}
	for {
		p := pair{}
		n, err := fmt.Fscanf(r, "%d-%d,%d-%d", &p.p1Start, &p.p1Stop, &p.p2Start, &p.p2Stop)

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}

		ret = append(ret, p)
	}

	return ret
}
