package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type op int8

const (
	ADD = iota
	SUB
	MULT
	DIV
)

type monkey struct {
	id        string
	op        op
	leftname  string //useful in first pass
	rightname string
	left      *monkey
	right     *monkey
	answer    *int //when we already have simple answer, no need to check left/right
}

func main() {

	in := "input.txt"
	file, _ := os.Open(in)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	ms := ParseMonkeys(lines)
	fmt.Println("Part1", ms["root"].Answer())

	fmt.Println("Part2 Intitial", ms["root"].left.Answer(), ms["root"].right.Answer())
	fmt.Println("Part2 humn", JiggleHumnTillMatch(ms))
}

func ParseMonkeys(lines []string) map[string]*monkey {
	ret := make(map[string]*monkey)
	ms := make([]monkey, len(lines))
	for i, line := range lines {
		split := strings.Split(line, ": ")
		ms[i].id = split[0]
		if ans, err := strconv.Atoi(split[1]); err == nil {
			ms[i].answer = &ans
		} else {
			split = strings.Split(split[1], " ")
			ms[i].leftname = split[0]
			ms[i].rightname = split[2]
			switch split[1] {
			case "+":
				ms[i].op = ADD
			case "-":
				ms[i].op = SUB
			case "*":
				ms[i].op = MULT
			case "/":
				ms[i].op = DIV
			default:
				panic("invalid op")
			}
		}
		ret[ms[i].id] = &ms[i]
	}

	for _, v := range ret {
		if v.answer == nil {
			v.left = ret[v.leftname]
			v.right = ret[v.rightname]
		}
	}
	return ret
}

func (m *monkey) Answer() int {
	if m.answer == nil {
		var ans int
		switch m.op {
		case ADD:
			ans = m.left.Answer() + m.right.Answer()
		case SUB:
			ans = m.left.Answer() - m.right.Answer()
		case MULT:
			ans = m.left.Answer() * m.right.Answer()
		case DIV:
			ans = m.left.Answer() / m.right.Answer()
		}
		// m.answer = &ans
		return ans
	}

	return *m.answer
}

func JiggleHumnTillMatch(ms map[string]*monkey) int {
	for x := -2; x < 2; x++ {
		v := x
		ms["humn"].answer = &v
		l, r := ms["root"].left.Answer(), ms["root"].right.Answer()
		fmt.Println("Jiggle", x, l == r, "diff", l-r, l, r)
		if l == r {
			return x
		}
	}

	//NOTE This Start value was found with some manual iteration with larger step sizes
	var prevDiff int = math.MaxInt64
	for x := 3243420789000; x < 1e13; x += 1 {
		v := x
		ms["humn"].answer = &v
		l, r := ms["root"].left.Answer(), ms["root"].right.Answer()
		diff := l - r

		if l == r {
			return x
		}

		if abs(diff) < abs(prevDiff) && x%1e3 == 0 {
			fmt.Printf("improved %d %e (prev %e)\n", x, float64(diff), float64(prevDiff))
		} else if abs(diff) > abs(prevDiff) {
			fmt.Println("Bad dir")
			return -1
		}
		prevDiff = diff
	}

	return -1
}
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
