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

	"github.com/samber/lo"
)

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	lines := parseInput(f)

	t := buildTree(lines)

	dirTotals := make(map[string]int)
	for k := range t {
		tallyDirs(dirTotals, t, filepath.Dir(k))
	}

	sum := lo.Sum(lo.Values(lo.PickBy(dirTotals, func(k string, v int) bool { return v < 100000 })))
	fmt.Println("part 1", sum)

	unused := 70000000 - lo.Sum(lo.Values(t))
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

func tallyDirs(tally map[string]int, files map[string]int, dirPath string) {
	total := lo.Sum(lo.Values(lo.PickBy(files, func(k string, v int) bool {
		return strings.HasPrefix(k, dirPath)
	})))
	tally[dirPath] = total
	parent := filepath.Dir(dirPath)
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

func parseInput(r io.Reader) []string {
	lines := []string{}
	s := bufio.NewScanner(r)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		lines = append(lines, s.Text())
	}

	return lines
}
