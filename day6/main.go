package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

var (
	part1Tests = []string{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb",
		"bvwbjplbgvbhsrlpgdmjqwftvncz",
		"nppdvjthqldpwncqszvftbrmjlhg",
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg",
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw",
	}
	part2Tests = []string{
		"mjqjpqmgbljsphdztnvjfqwrcgsmlb",
		"bvwbjplbgvbhsrlpgdmjqwftvncz",
		"nppdvjthqldpwncqszvftbrmjlhg",
		"nznrnfrfntjfmvfwmzdfjlvtqnbhcprsg",
		"zcfzfwzzqfrljwzlrfnpqdbhtmscgvjw",
	}
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	lines := parseInput(f)

	for _, l := range part1Tests {
		fmt.Println("Part1 start", findStart(l))
	}
	fmt.Println("Part1 Solution", findStart(lines[0]))

	for _, l := range part2Tests {
		fmt.Println("Part2 start", findStartPart2(l))
	}
	fmt.Println("Part2 Solution", findStartPart2(lines[0]))

}

func findStart(line string) int {
	for i := 4; i < len(line); i++ {
		if countUniqueChars(line[i-4:i]) == 4 {
			return i
		}
	}
	return -1
}

func findStartPart2(line string) int {
	for i := 14; i < len(line); i++ {
		if countUniqueChars(line[i-14:i]) == 14 {
			return i
		}
	}
	return -1
}
func countUniqueChars(str string) int {
	m := make(map[rune]struct{})
	for _, c := range str {
		m[c] = struct{}{}
	}
	return len(m)
}

func parseInput(r io.Reader) []string {
	lines := []string{}
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return lines
}
