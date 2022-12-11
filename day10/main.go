package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type cycleList []int

// Get the accumulated value for X during the nth cycle
func (c cycleList) GetCycleTotal(n int) int {
	if n <= 1 {
		return 1
	}
	return c[n-2]
}

func (c cycleList) Print() {
	for i := 0; i < len(c); i++ {
		si := i % 40

		pos := c.GetCycleTotal(i + 1)
		delta := pos - si
		if delta > 1 || delta < -1 {
			fmt.Print(".")
		} else {
			fmt.Print("#")
		}
		if si == 39 {
			fmt.Println()
		}
	}
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	lines := parseInput(f)
	cycles := buildCycles(lines)

	total := 0
	for _, i := range []int{20, 60, 100, 140, 180, 220} {
		xi := cycles.GetCycleTotal(i)
		sig := i * xi
		total += sig
		fmt.Println("Getting During Cycle", i, xi, "sig", sig, total)
	}

	cycles.Print()
}

func buildCycles(lines []string) cycleList {
	total := 1 //start with 1
	ret := make(cycleList, 0)
	for _, line := range lines {
		var x int
		if line == "noop" {
			ret = append(ret, total) //noop doesn't add anything to x
		} else {
			if n, err := fmt.Sscanf(line, "addx %d", &x); err != nil || n != 1 {
				log.Fatal("Failed to parse input")
			}
			ret = append(ret, total)
			total += x
			ret = append(ret, total)
		}
	}
	return ret
}

func parseInput(f io.ReadSeekCloser) []string {
	lines := []string{}
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		lines = append(lines, s.Text())
	}
	return lines
}
