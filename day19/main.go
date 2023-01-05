package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

type blueprint struct {
	id                      int
	OreRobotCost            int
	ClayRobotCost           int
	ObsidianRobotCostsOre   int
	ObsidianRobotCostsClay  int
	GeodeRobotCostsOre      int
	GeodeRobotCostsObsidian int
}

type state struct {
	step int

	oreRobot      int
	clayRobot     int
	obsidianRobot int
	geodeRobot    int

	dead *int

	ore      int
	clay     int
	obsidian int
	geode    int
}

func main() {
	in := "input.txt"
	file, _ := os.Open(in)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	bps := make([]blueprint, 0)
	for scanner.Scan() {
		bp := blueprint{}
		if n, err := fmt.Sscanf(scanner.Text(), `Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.`,
			&bp.id, &bp.OreRobotCost, &bp.ClayRobotCost, &bp.ObsidianRobotCostsOre, &bp.ObsidianRobotCostsClay, &bp.GeodeRobotCostsOre, &bp.GeodeRobotCostsObsidian); err != nil || n != 7 {
			log.Fatal("parse fail")
		}
		bps = append(bps, bp)
	}

	w := sync.WaitGroup{}
	w.Add(len(bps))
	sum := 0
	for _, bp := range bps {
		go func(bfunc blueprint) {
			s := bfunc.MaxGeodes(24)
			sum += s * bfunc.id
			fmt.Println("blueprint", bfunc, "best", s)
			w.Done()
		}(bp)
	}
	w.Wait()
	fmt.Println("Part1 Sum", sum)
}

func (b blueprint) MaxGeodes(steps int) int {
	s := state{oreRobot: 1}

	return s.nextGeodes(&b, steps)
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func (s *state) MinutePassesCollect() {
	s.ore += s.oreRobot
	s.clay += s.clayRobot
	s.obsidian += s.obsidianRobot
	s.geode += s.geodeRobot
	s.step++
}

// param is a pointer cause this is a hot path and less memory alloc
func (s state) nextGeodes(b *blueprint, maxSteps int) int {
	if s.step >= maxSteps {
		return s.geode
	}

	ret := s.geode

	//DESIGN REVIEW - Greedy here might build all Obsidian robots (i.e. if we need more ore for geode robot than obsidian)
	switch {
	case s.ore >= b.GeodeRobotCostsOre && s.obsidian >= b.GeodeRobotCostsObsidian:
		buildGeoRobot := s
		buildGeoRobot.ore -= b.GeodeRobotCostsOre
		buildGeoRobot.obsidian -= b.GeodeRobotCostsObsidian
		buildGeoRobot.MinutePassesCollect()
		buildGeoRobot.geodeRobot++

		// fmt.Println("Build Geo Robot", buildGeoRobot)

		ret = max(ret, buildGeoRobot.nextGeodes(b, maxSteps))

	case s.ore >= b.ObsidianRobotCostsOre && s.clay >= b.ObsidianRobotCostsClay:
		// &&
		//(s.obsidianRobot == 0 || b.
		buildObsRobot := s
		buildObsRobot.ore -= b.ObsidianRobotCostsOre
		buildObsRobot.clay -= b.ObsidianRobotCostsClay
		buildObsRobot.MinutePassesCollect()
		buildObsRobot.obsidianRobot++
		ret = max(ret, buildObsRobot.nextGeodes(b, maxSteps))

	default:
		if s.ore >= b.OreRobotCost {
			buildOreRobot := s
			buildOreRobot.ore -= b.OreRobotCost
			buildOreRobot.MinutePassesCollect()
			buildOreRobot.oreRobot++
			ret = max(ret, buildOreRobot.nextGeodes(b, maxSteps))
		}

		if s.ore >= b.ClayRobotCost {
			buildClayRobot := s
			buildClayRobot.ore -= b.ClayRobotCost
			buildClayRobot.MinutePassesCollect()
			buildClayRobot.clayRobot++
			ret = max(ret, buildClayRobot.nextGeodes(b, maxSteps))
		}
	}
	{
		cantBuildAny := s
		cantBuildAny.MinutePassesCollect()

		ret = max(ret, cantBuildAny.nextGeodes(b, maxSteps))
	}

	return ret
}
