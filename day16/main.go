package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/samber/lo"
)

type graph []*node

type node struct {
	valveName string
	flowRate  int
	tunnels   []string
	edges     []*edge
}

type edge struct {
	cost int
	src  *node
	dest *node
}

type move struct {
	here *node
	prev *node
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}
	g := parseInput(f)
	fmt.Println("nodes", g)

	fmt.Println("Best", g.WalkAll())
}

func (g graph) Get(valveName string) *node {
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

var bestYet int = 0

func (g graph) WalkAll() int {
	best := 0

	// visited := make(map[string]bool)

	minutes := 30
	t0 := time.Now()

	start := g.Get("AA")

	g.recurse("", minutes, 0, move{start, start}, move{start, start})

	// g.Example(minutes, 0, "", bestMoves...)

	//try open n, try move through all tunnels that we've not aleady visited
	// visited[n.valveName] = true

	// for _, t := range n.tunnels {
	// 	best = max(best, (minutes-1)*n.flowRate + g.recurse(minutes-2, g.Get(t))
	// 	best = max(best, g.recurse(minutes-1, g.Get(t)) //TODO: Come back to valve and open if enough time??
	// }

	fmt.Println("Runtime: ", time.Since(t0))
	fmt.Println("Best: ", lo.Max(lo.Values(ALL)))
	fmt.Println("Best int: ", bestYet)

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

func (g graph) recurse(valvesOpened string, minutes int, score int, m, m2 move) {
	if maxMoves++; maxMoves > 2000000000 {
		fmt.Println("ABORT")
		return
	}

	if minutes < 2 {
		// key := visited
		// ALL[key] = score
		bestYet = max(bestYet, score)
		// fmt.Println(pre, minutes, "Final Score ", score, key)
		return
	}

	here := m.here
	if here.flowRate > 0 && strings.Index(valvesOpened, "+"+here.valveName) == -1 {
		n := move{prev: here, here: here}

		// for _, e2 := range m2.here.edges {
		// 	//Don't about-face if we didn't even open the here valve
		// 	if m2.prev == e2.dest {
		// 		continue
		// 	}

		// 	n2 := move{prev: m2.here, here: e2.dest}

		n2 := m2 //FIXME

		g.recurse(valvesOpened+">+"+here.valveName, minutes-1, score+((minutes-1)*here.flowRate), n, n2)
		// }

		// fmt.Println(pre, minutes, "Got score ", rc0, "for opening valve", here.valveName, "before moving to", t, visited)
	}

	for _, e := range m.here.edges {
		//Don't about-face if we didn't even open the here valve
		if m.prev == e.dest {
			continue
		}

		n := move{prev: m.here, here: e.dest}
		// for _, e2 := range m2.here.edges {
		// 	//Don't about-face if we didn't even open the here valve
		// 	if m2.prev == e2.dest {
		// 		continue
		// 	}

		// 	n2 := move{prev: m2.here, here: e2.dest}

		n2 := m2 //FIXME

		g.recurse(valvesOpened, minutes-1, score, n, n2)
		// fmt.Println(pre, minutes, "Got score ", rc1, "for skipping valve", here.valveName, "before moving to", t, visited)
		// }
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
	for s.Scan() {

		line := s.Text()
		n := node{}
		if i, err := fmt.Sscanf(line, "Valve %s has flow rate=%d", &n.valveName, &n.flowRate); i < 2 || err != nil {
			log.Fatal("Parse fail - only got ", i, "values - err:", err)
		}

		rightHalf := strings.Split(line, "; ")[1]
		splits := strings.Split(rightHalf, " ")

		valves := strings.Join(splits[4:], "")

		n.tunnels = strings.Split(valves, ",")
		ret = append(ret, &n)
	}

	//Help the greedy
	// for _, v := range ret {
	// 	v.SortTunnels(ret)
	// }

	for _, n := range ret {
		fmt.Println("Valve ", n.valveName, "flow", n.flowRate, "tunnels", n.tunnels)
		for _, t := range n.tunnels {
			dest := ret.Get(t)
			e := ret.collapseEdge(n, dest)
			fmt.Println("Edge: ", e.src.valveName, e.dest.valveName, e.cost)

			n.edges = append(n.edges, e)
		}
	}
	return ret
}

func (g graph) collapseEdge(src *node, dest *node) *edge {
	/*if false && dest.flowRate == 0 && len(dest.tunnels) == 2 {
		var next *node
		for i := 0; i < len(dest.tunnels); i++ {
			if dest.tunnels[i] != src.valveName {
				next = g.Get(dest.tunnels[i])
			}
		}
		e := g.collapseEdge(dest, next)
		ret := &edge{src: src,
			dest: e.dest,
			cost: e.cost + 1,
		}
		// fmt.Println("Collapsed Edge: ", e.src.valveName, e.dest.valveName, ret.cost)
		return ret
	} else */{
		return &edge{
			src:  src,
			dest: dest,
			cost: 1,
		}
	}
}
