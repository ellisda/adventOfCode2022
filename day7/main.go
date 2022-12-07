package main

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	lines := parseInput(f)

	t := buildTree(lines)
	fmt.Println("tree", t)

	dirTotals := make(map[string]int)
	for k := range t {
		tallyDirs(dirTotals, t, filepath.Dir(k))
	}
	sum := 0
	for k, v := range dirTotals {
		if v < 100000 {
			sum += v
			fmt.Println("Summing", v, k)
		}

	}
	fmt.Println("part 1", sum)

	unused := 70000000 - sumDir(t, "/")
	needed := 30000000 - unused
	fmt.Println("Part 2 Need", needed)
	smallestAcceptable := math.MaxInt
	for k, v := range dirTotals {
		if v >= needed && v < smallestAcceptable {
			smallestAcceptable = v
			fmt.Println("Best Yet", v, k)
		}
	}
}

func tallyDirs(tally map[string]int, files map[string]int, leaf string) {
	total := sumDir(files, leaf)
	tally[leaf] = total
	parent := filepath.Dir(leaf)
	if len(parent) > 2 {
		tallyDirs(tally, files, parent)
	}
}

func buildTree(lines []string) map[string]int {
	cwd := ""
	ret := make(map[string]int)
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "$ cd "):
			param := line[5:]
			if param == ".." {
				cwd = path.Dir(cwd)
			} else {
				cwd = path.Join(cwd, param)
			}
		case line == "$ ls":
			continue
		case strings.HasPrefix(line, "dir "):
			continue

		default:
			size, name := 0, ""
			fmt.Sscanf(line, "%d %s", &size, &name)
			full := filepath.Join(cwd, name)
			fmt.Println("Got:", full, size, "--", line)
			ret[full] = size
		}

	}
	return ret
}

func sumDir(tree map[string]int, dirPath string) int {
	total := 0
	for k, v := range tree {
		if strings.HasPrefix(k, dirPath) {
			total += v
		}
	}
	fmt.Println("Sum", dirPath, total)
	return total
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
