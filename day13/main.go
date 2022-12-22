package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

const (
	OUTOFORDER = -1
	INORDER    = 1
)

type node struct {
	value    int
	children []*node
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	fmt.Println("test", "[[],[[],3,[]],[]]")
	fmt.Println("unroll", ParseTree2("[[],[[],3,[]],[]]").print())

	fmt.Println("test", "[[[20],[[3,1,6,8,1],9,9,[4,7]],9,[]],[[[3]],[9,3,[1,8,7,6,4],[0,1,1]],5]]")
	fmt.Println("test", ParseTree2("[[[20],[[3,1,6,8,1],9,9,[4,7]],9,[]],[[[3]],[9,3,[1,8,7,6,4],[0,1,1]],5]]").print())

	left := ParseTree2("[[[],4,[4],[[4,0,8,5,3],[10],8]],[0,1,9,[[6,8,9,4],6,[]]],[[[5,4,9,0,9],[2],2,2],[[1,9,10,7,7],[3,7,5,1],[6,5,0,0],4,1],0]]")
	right := ParseTree2("[[[],1,10,3],[[[8],10,5,[2,10,7,3],0]],[7,[[0,1,8,10],1,1,9,[]]],[[[0],6],[[]],5]]")
	fmt.Println("expected false, actual", treeCompare(left, right))

	lines := parseInput(f)

	part1 := inOrderPairs(lines)
	fmt.Println(part1)
	sum := 0
	for i, b := range part1 {
		if b {
			sum += i + 1
		}
	}
	fmt.Println("part1 sum", sum)

	part2 := 1
	sorted := sortPackets(lines, "[[2]]", "[[6]]")
	for i, x := range sorted {
		fmt.Println(i, x.print())
		switch x.print() {
		case "[[2]]":
			fallthrough
		case "[[6]]":
			fmt.Println("found divisor at", i, x.print())
			part2 *= (i + 1)
		}
	}

	fmt.Println("Part2", part2)
}

func (n *node) print() string {
	if n.children == nil {
		return fmt.Sprintf("%d", n.value)
	} else {
		ret := "["
		for _, n2 := range n.children {
			ret += n2.print() + ","
		}

		return strings.TrimRight(ret, ",") + "]"
	}
}

func ParseTree2(s string) *node {
	stack := []*node{}
	var cur *node = &node{-2, nil}
	var span string
	for _, c := range s {
		switch c {
		case '[':
			for _, v := range split(span) {
				cur.children = append(cur.children, &node{v, nil})
			}
			span = ""

			newChild := &node{-1, []*node{}}
			cur.children = append(cur.children, newChild)
			stack = append(stack, cur)
			cur = newChild
		case ']':
			for _, v := range split(span) {
				cur.children = append(cur.children, &node{v, nil})
			}
			span = ""

			cur = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		default:
			span = span + string(c)
		}

	}

	cur = cur.children[0]
	// fmt.Println("Inner", cur.print())
	// fmt.Println()
	return cur
}

func split(s string) []int {
	ss := strings.Split(strings.Trim(s, ","), ",")
	if len(ss) == 0 || len(ss) == 1 && ss[0] == "" {
		return []int{}
	}
	return lo.Map(ss, conv)
}

func unroll(n string) []string {
	// firstSub := strings.Split()

	if strings.HasPrefix(string(n), "[") {
		return strings.Split(n[1:len(n)-1], ",")
	}
	return strings.Split(n, ",")
}

func inOrderPairs(lines []string) []bool {
	numPairs := (len(lines) + 1) / 3
	ret := make([]bool, numPairs)
	for i := 0; i < numPairs; i++ {
		l := ParseTree2(lines[3*i])
		r := ParseTree2(lines[3*i+1])

		ret[i] = treeCompare(l, r) != OUTOFORDER

		// ret[i] = inOrderPair(lines[3*i], lines[3*i+1])
	}
	return ret
}

func treeCompare(left, right *node) int {
	fmt.Println("Left", left.print())
	fmt.Println("right", right.print())

	if left.children == nil && right.children == nil {
		if left.value < right.value {
			fmt.Println("true")
			return INORDER
		} else if left.value > right.value {
			fmt.Println("false")
			return OUTOFORDER
		} else {
			return 0
		}
	}
	if (left.children == nil) != (right.children == nil) {
		if left.children == nil {
			return treeCompare(&node{-1, []*node{left}}, right)
		} else {
			return treeCompare(left, &node{-1, []*node{right}})
		}
	}

	for i := 0; i < len(left.children); i++ {

		if i >= len(right.children) {
			fmt.Println("false - ran out of right side at index", i)
			return OUTOFORDER
		}

		temp := treeCompare(left.children[i], right.children[i])
		if temp != 0 {
			return temp
		}
	}

	if len(left.children) < len(right.children) {
		fmt.Println("true ran out of left list", len(left.children), len(right.children))
		return INORDER
	}

	fmt.Println("neutral len", len(left.children), len(right.children))
	return 0

}

func sortPackets(lines []string, dividers ...string) []*node {
	numPairs := (len(lines) + 1) / 3
	packets := make([]*node, 0, numPairs*2)
	for i := 0; i < numPairs; i++ {
		packets = append(packets, ParseTree2(lines[3*i]))
		packets = append(packets, ParseTree2(lines[3*i+1]))
	}
	for _, d := range dividers {
		packets = append(packets, ParseTree2(d))
	}
	// ret[i] = treeCompare(l, r) != OUTOFORDER

	sort.Slice(packets, func(i, j int) bool { return treeCompare(packets[i], packets[j]) == INORDER })
	return packets
}

func conv(in string, i int) int {
	v, _ := strconv.Atoi(in)
	return v
}

func trim(in string) string {
	return strings.ReplaceAll(strings.ReplaceAll(in, "[", ""), "]", "")
}

func inOrderPair(left, right string) bool {
	fmt.Println("Left", left, "right", right)

	l := strings.Split(trim(left), ",")
	r := strings.Split(trim(right), ",")

	fmt.Println("Left", l, "right", r)
	l0 := lo.Map(l, conv)
	r0 := lo.Map(r, conv)
	fmt.Println("Left", l0, "right", r0)
	for i := 0; i < len(l0); i++ {
		if i >= len(r0) {
			fmt.Println("false - ran out of right side at index", i)
			return false
		}

		if l0[i] > r0[i] {
			fmt.Println("false")
			return false
		}
		if l0[i] < r0[i] {
			fmt.Println("true")
			return true
		}
	}

	if len(l) >= len(r) {
		fmt.Println("false len", len(l), len(r))
		return false
	}
	fmt.Println("true len", len(l), len(r))
	return true
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
