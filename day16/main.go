package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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

func (n *node) SortTunnels(g graph) {
	sort.Slice(n.tunnels, func(i, j int) bool {
		ni := g.Get(n.tunnels[i])
		nj := g.Get(n.tunnels[j])
		return ni.flowRate > nj.flowRate
	})
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

var bestMoves = []string{
	"DD", "+DD", "CC", "BB", "+BB", "AA", "II", "JJ", "+JJ",
	"II", "AA", "DD", "EE", "FF", "GG", "HH", "+HH",
	"GG", "FF", "EE", "+EE", "DD", "CC", "+CC",
	// "DD", "CC", "DD", "CC", "DD", "CC", "DD", "CC", //who cares about the last 6-8 here
}

func (g graph) WalkAll() int {
	best := 0

	// visited := make(map[string]bool)

	minutes := 30

	g.recurse("", ">"+g[0].valveName, minutes, 0, g[0].valveName)

	// g.Example(minutes, 0, "", bestMoves...)

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

func (g graph) recurse(pre string, visited string, minutes int, score int, at string) {
	if maxMoves++; maxMoves > 900000000 {
		return
	}

	if /*has(visited, "+"+here.valveName) ||*/ minutes < 2 {
		key := visited
		ALL[key] = score
		// if score >= 1649 {

		// fmt.Println(pre, minutes, "Final Score ", score, key)
		// }
		return
	}

	if visited == ">AA>DD>+DD>CC>BB" {
		fmt.Println("On a good path (moves: ", maxMoves, ")", visited)
	}
	if visited == ">AA>DD>+DD>CC>BB>+BB>AA>II>JJ>+JJ" {
		fmt.Println("On a good path (moves: ", maxMoves, ")", visited)
	}
	if visited == ">AA>DD>+DD>CC>BB>+BB>AA>II>JJ>+JJ>II" {
		fmt.Println("On a good path (moves: ", maxMoves, ")", visited)
	}
	if visited == ">AA>DD>+DD>CC>BB>+BB>AA>II>JJ>+JJ>II>AA" {
		fmt.Println("On a good path (moves: ", maxMoves, ")", visited)
	}
	if visited == ">AA>DD>+DD>CC>BB>+BB>AA>II>JJ>+JJ>II>AA>DD" {
		fmt.Println("On a good path (moves: ", maxMoves, ")", visited)
	}
	if visited == ">AA>DD>+DD>CC>BB>+BB>AA>II>JJ>+JJ>II>AA>DD>EE>FF>GG>HH>+HH" {
		fmt.Println("On a good path (moves: ", maxMoves, ")", visited)
	}

	// best := 0
	// visited = visited + ">" + here.valveName

	// var here node
	// if strings.HasPrefix(at, "+") {
	// 	here = g.Get(at[1:])

	// 	// minutes--
	// 	// score += minutes * here.flowRate
	// 	// // visited = visited+">+"+here.valveName
	// 	// g.recurse(pre+".", visited+">"+at, minutes-1, score, at[1:])
	// } else {
	// }
	here := g.Get(at)

	if here.flowRate > 0 && strings.Index(visited, "+"+here.valveName) == -1 {
		g.recurse(pre+".", visited+">+"+here.valveName, minutes-1, score+((minutes-1)*here.flowRate), here.valveName)
		// fmt.Println(pre, minutes, "Got score ", rc0, "for opening valve", here.valveName, "before moving to", t, visited)
	}

	for _, t := range here.tunnels {
		//FIXME - Can't avoid walking an egde twice, have to allow ">AA>DD>...>+JJ>..>AA>DD"
		// nextHop := here.valveName + ">" + t
		// if strings.Index(visited, nextHop) != -1 {
		// 	continue
		// }
		//FIXME - Avoiding backtracking at all is wrong, we need to allow ">II>JJ>+JJ>II"
		// reverseHop := ">" + t + ">" + here.valveName
		// if strings.Index(visited, reverseHop) != -1 {
		// 	continue
		// }
		//Don't about-face if we didn't even open the here valve
		reverseHop := ">" + t + ">" + here.valveName
		if strings.HasSuffix(visited, reverseHop) {
			continue
		}

		// next := g.Get(t)
		// if t == "DD" {
		// 	fmt.Print("HERE")
		// }
		// if len(visited) > 0 && strings.HasSuffix(visited, t+">"+here.valveName) {
		// 	continue //Don't turn 180 and go back where you came from
		// }

		g.recurse(pre+".", visited+">"+t, minutes-1, score, t)
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

	//Help the greedy
	for _, v := range ret {
		(&v).SortTunnels(ret)
	}

	return ret
}
