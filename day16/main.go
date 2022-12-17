package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/samber/lo"
)

type graph []node

type node struct {
	valveName string
	flowRate  int
	tunnels   []string
}

func main() {
	f, err := os.Open("input2.txt")
	if err != nil {
		panic(err)
	}
	g := parseInput(f)
	fmt.Println("nodes", g)

	fmt.Println("Best", g.WalkAll())
}

func (g graph) Get(valveName string) node {
	for _, n := range g {
		if n.valveName == valveName {
			return n
		}
	}
	panic("Can't find node" + valveName)
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func (g graph) WalkAll() int {
	best := 0

	// visited := make(map[string]bool)

	minutes := 30

	// g.recurse("", "", minutes, 0, g[0])

	g.Example(minutes, 0, "", "DD", "+DD", "CC", "BB", "+BB", "AA", "II", "JJ", "+JJ",
		"II", "AA", "DD", "EE", "FF", "GG", "HH", "+HH",
		"GG", "FF", "EE", "+EE", "DD", "CC", "+CC",
		"DD", "CC", "DD", "CC", "DD", "CC", "DD", "CC")

	//try open n, try move through all tunnels that we've not aleady visited
	// visited[n.valveName] = true

	// for _, t := range n.tunnels {
	// 	best = max(best, (minutes-1)*n.flowRate + g.recurse(minutes-2, g.Get(t))
	// 	best = max(best, g.recurse(minutes-1, g.Get(t)) //TODO: Come back to valve and open if enough time??
	// }

	fmt.Println("Best: ", lo.Max(lo.Values(ALL)))

	return best

}

func has(list string, s string) bool {
	for _, l := range strings.Split(list, ">") {
		if s == l {
			return true
		}
	}
	return false
}

var ALL map[string]int = make(map[string]int, 0)
var maxMoves int

func (g graph) Example(minutes int, score int, visited string, moves ...string) {

	if len(moves) == 0 || minutes < 2 {
		key := visited //+ ">" + here.valveName
		ALL[key] = score
		fmt.Println(minutes, "Final Score ", score, key)
		return
	}

	m := moves[0]
	if strings.HasPrefix(m, "+") {
		next := g.Get(m[1:])
		g.Example(minutes-1, score+((minutes-1)*next.flowRate), visited+">"+m, moves[1:]...)
		// 	g.recurse(pre+".", visited+">+"+here.valveName, minutes-2, score+(minutes-1)*here.flowRate, next)
		// 	// fmt.Println(pre, minutes, "Got score ", rc0, "for opening valve", here.valveName, "before moving to", t, visited)
	} else {

		g.Example(minutes-1, score, visited+">"+m, moves[1:]...)

		// g.recurse(pre+".", visited+">"+here.valveName, minutes-1, score, next)
	}

}

func (g graph) recurse(pre string, visited string, minutes int, score int, here node) {
	if maxMoves++; maxMoves > 900000000 {
		return
	}

	if /*has(visited, "+"+here.valveName) ||*/ minutes < 2 {
		key := visited + ">" + here.valveName
		ALL[key] = score
		if score >= 1649 {

			fmt.Println(pre, minutes, "Final Score ", score, key)
		}
		return
	}

	// best := 0
	// visited = visited + ">" + here.valveName

	for _, t := range here.tunnels {
		nextHop := here.valveName + ">" + t
		if strings.Index(visited, nextHop) != -1 {
			continue
		}

		next := g.Get(t)
		// if t == "DD" {
		// 	fmt.Print("HERE")
		// }
		// if len(visited) > 0 && strings.HasSuffix(visited, t+">"+here.valveName) {
		// 	continue //Don't turn 180 and go back where you came from
		// }

		if strings.Index(visited, "+"+here.valveName) == -1 {
			g.recurse(pre+".", visited+">+"+here.valveName, minutes-2, score+(minutes-1)*here.flowRate, next)
			// fmt.Println(pre, minutes, "Got score ", rc0, "for opening valve", here.valveName, "before moving to", t, visited)
		}

		g.recurse(pre+".", visited+">"+here.valveName, minutes-1, score, next)
		// fmt.Println(pre, minutes, "Got score ", rc1, "for skipping valve", here.valveName, "before moving to", t, visited)
		// best = max(rc1, best)

	}
}

// func (g graph) recurse(pre string, visited map[string]bool, minutes int, here node) int {
// 	if visited[here.valveName] {
// 		return 0
// 	}

// 	best := 0
// 	visited[here.valveName] = true

// 	for _, t := range here.tunnels {
// 		next := g.Get(t)
// 		if minutes >= 2 {
// 			rc0 := (minutes-1)*here.flowRate + g.recurse(pre+".", visited, minutes-2, next)
// 			fmt.Println(pre, minutes, "Got score ", rc0, "for opening valve", here.valveName, "before moving to", t)
// 			best = max(rc0, best)
// 		}
// 		if minutes >= 3 {
// 			rc1 := g.recurse(pre+".", visited, minutes-1, next) //TODO: Come back to valve and open if enough time??
// 			fmt.Println(pre, minutes, "Got score ", rc1, "for skipping valve", here.valveName, "before moving to", t)
// 			best = max(rc1, best)
// 		}
// 	}
// 	return best
// }

func parseInput(f io.ReadSeekCloser) graph {
	ret := graph{}
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	n := node{}
	for s.Scan() {

		line := s.Text()
		if i, err := fmt.Sscanf(line, "Valve %s has flow rate=%d", &n.valveName, &n.flowRate); i < 2 || err != nil {
			log.Fatal("Parse fail - only got ", i, "values - err:", err)
		}

		rightHalf := strings.Split(line, "; ")[1]
		splits := strings.Split(rightHalf, " ")

		valves := strings.Join(splits[4:], "")

		n.tunnels = strings.Split(valves, ",")
		ret = append(ret, n)
	}
	return ret
}
