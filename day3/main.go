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

	rucksacks := parseInput(f)
	fmt.Println("found: ", rucksacks)

	total := 0
	for _, r := range rucksacks {
		dups := r.Duplicates()
		fmt.Println("From", r, " we find Duplicates", dups)
		for _, d := range dups {
			total += d
		}
	}
	fmt.Println("Total", total)

	groups := make([]group, len(rucksacks)/3)
	for i, r := range rucksacks {
		groups[i/3][i%3] = r
	}

	fmt.Println(groups)
	totalBadgeScore := 0
	for _, g := range groups {
		badge := g.Badge()
		fmt.Println("Badge", string(badge), g)
		totalBadgeScore += runeScore(badge)
	}
	fmt.Println("Total Badge Score", totalBadgeScore)
}

type group [3]rucksack

func (g group) Badge() rune {
	m1, m2 := make(map[rune]struct{}), make(map[rune]struct{})
	for _, r := range g[0].left + g[0].right {
		m1[r] = struct{}{}
	}
	for _, r := range g[1].left + g[1].right {
		m2[r] = struct{}{}
	}
	for _, r := range g[2].left + g[2].right {
		if _, ok1 := m1[r]; !ok1 {
			continue
		}
		if _, ok2 := m2[r]; !ok2 {
			continue
		}
		return r
	}
	panic("missing badge")
}

type rucksack struct {
	left  string
	right string
}

func (r rucksack) Duplicates() []int {
	seen := make(map[rune]struct{})
	dups := make(map[rune]struct{})
	ret := []int{}
	for _, r := range r.left {
		seen[r] = struct{}{}
	}

	for _, r := range r.right {
		if _, ok := seen[r]; ok {
			dups[r] = struct{}{}
		}
	}
	for k, _ := range dups {
		ret = append(ret, runeScore(k))
	}

	return ret
}

func runeScore(r rune) int {
	ret := int(r) - int('a')
	switch {
	case int(r) >= int('a') && int(r) <= int('z'):
		ret = 1 + int(r) - int('a')
	case int(r) >= int('A') && int(r) <= int('Z'):
		ret = 27 + int(r) - int('A')

	}
	fmt.Println("rune score", string(r), int(r), ret)
	return ret
}

func parseInput(r io.Reader) []rucksack {
	ret := []rucksack{}
	for {
		var s string

		n, err := fmt.Fscanln(r, &s)

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}

		l, r := s[:len(s)/2], s[len(s)/2:]
		ret = append(ret, rucksack{l, r})
	}

	return ret
}
