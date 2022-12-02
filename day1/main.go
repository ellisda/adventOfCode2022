package main

import (
	"fmt"
	"io"
	"os"
	"sort"
)

func main() {

	fmt.Println("Hello")

	f, err := os.Open("input.txt")
	if err != nil {

	}

	elves := parseLineGroups(f)

	all := make([]int, len(elves))
	max := 0
	for i, e := range elves {
		sum := e.Sum()
		all[i] = sum

		if sum > max {
			max = sum
		}
	}
	fmt.Println("Biggest ElfCals: ", max)

	sort.Slice(all, func(x, y int) bool { return all[x] > all[y] })
	fmt.Printf("Biggest 3 Elves (sum=%d): %d %d %d\n", all[0]+all[1]+all[2], all[0], all[1], all[2])
}

type elfCals []int

func (e elfCals) Sum() int {
	ret := 0
	for _, c := range e {
		ret += c
	}
	return ret
}

func parseLineGroups(f io.Reader) []elfCals {
	ret := []elfCals{}
	ret = append(ret, elfCals{})
	groups := 0
	for {
		var i int
		n, err := fmt.Fscanln(f, &i)

		if err != nil && err.Error() == "unexpected newline" {
			groups++
			ret = append(ret, elfCals{})
			continue
		}

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}

		ret[groups] = append(ret[groups], i)

	}

	return ret
}
