package main

type predicate func(from, to *square) bool

func (g grid) dfs(end *square, pFunc predicate) map[position]int {
	scores := make(map[position]int)
	visited := make(map[position]bool)
	g.dfs_recurse(end, visited, scores, 0, pFunc)
	return scores
}

func (g grid) dfs_recurse(current *square, visited map[position]bool, scores map[position]int, score int, pFunc predicate) {
	if s, ok := scores[current.position()]; ok && s <= score {
		return // abandon this path, we've already seen a better score
	}

	visited[current.position()] = true
	scores[current.position()] = score

	for _, next := range getNextMoves(g, current) {
		if !pFunc(current, next) {
			g.dfs_recurse(next, visited, scores, score+1, pFunc)
		}
	}
}

func getNextMoves(g grid, from *square) []*square {
	ret := make([]*square, 0, 4)
	for _, p := range []position{
		{from.x + 1, from.y},
		{from.x - 1, from.y},
		{from.x, from.y - 1},
		{from.x, from.y + 1},
	} {
		if p.x < 0 || p.y < 0 || p.x >= len(g[0]) || p.y >= len(g) {
			continue
		}
		ret = append(ret, g[p.y][p.x])
	}
	return ret
}

// needsGear tells us whether we're unable to walk forward from one square to another without gear
func needsGear(from, to *square) bool {
	switch {
	case to.elevation == End:
		return from.elevation != 'z'
	case from.elevation == Start:
		return false
	default:
		return to.elevation > (from.elevation + 1)
	}
}

// When walking problem backwards, are we unable to proceed without gear
func needsGearBackward(from, to *square) bool {
	return needsGear(to, from)
	// switch {
	// case from.elevation == End:
	// 	return to.elevation == 'z'
	// case to.elevation == Start:
	// 	return true
	// case from.elevation - 1 <= to.elevation:
	// 	return true
	// }
}
