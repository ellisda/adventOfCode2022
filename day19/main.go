package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
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

var GLOBAL state

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

	t0 := time.Now()
	w := sync.WaitGroup{}
	w.Add(len(bps))
	sum := 0
	for _, bp := range bps {
		go func(bfunc blueprint) {
			s := bfunc.MaxGeodes(24)
			sum += s * bfunc.id
			fmt.Println("blueprint", bfunc, "max ore", max(max(max(bfunc.GeodeRobotCostsOre, bfunc.ObsidianRobotCostsOre), bfunc.ClayRobotCost), bfunc.OreRobotCost), "best", s)
			w.Done()
		}(bp)
	}
	w.Wait()
	fmt.Println("Part1 Sum", sum, "runtime", time.Since(t0))
	t1 := time.Now()

	// fmt.Println("Part2", state{oreRobot: 1}.nextGeodes_part2(&bps[0], 32))
	pt2 := bps[0:3]
	w.Add(len(pt2))
	mult := 1
	for _, bp := range pt2 {
		go func(bfunc blueprint) {
			s := bfunc.MaxGeodes(32)
			mult = mult * s
			fmt.Println("blueprint", bfunc, "max ore", max(max(max(bfunc.GeodeRobotCostsOre, bfunc.ObsidianRobotCostsOre), bfunc.ClayRobotCost), bfunc.OreRobotCost), "best", s)
			w.Done()
		}(bp)
	}
	w.Wait()
	fmt.Println("Part2", pt2, "mult", mult, "runtime", time.Since(t1))

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

	maxOrePerMin := max(max(max(b.GeodeRobotCostsOre, b.ObsidianRobotCostsOre), b.ClayRobotCost), b.OreRobotCost)
	maxClayPerMin := b.ObsidianRobotCostsClay
	maxObsPerMin := b.GeodeRobotCostsObsidian
	_ = maxClayPerMin + maxOrePerMin + maxObsPerMin

	//DESIGN REVIEW - Greedy here might build all Obsidian robots (i.e. if we need more ore for geode robot than obsidian)
	canBuildGeo := s.ore >= b.GeodeRobotCostsOre && s.obsidian >= b.GeodeRobotCostsObsidian
	canBuildObs := s.ore >= b.ObsidianRobotCostsOre && s.clay >= b.ObsidianRobotCostsClay
	canBuildOre := s.ore >= b.OreRobotCost
	canBuildClay := s.ore >= b.ClayRobotCost
	if canBuildGeo {
		buildGeoRobot := s
		buildGeoRobot.ore -= b.GeodeRobotCostsOre
		buildGeoRobot.obsidian -= b.GeodeRobotCostsObsidian
		buildGeoRobot.MinutePassesCollect()
		buildGeoRobot.geodeRobot++

		// fmt.Println("Build Geo Robot", buildGeoRobot)
		GLOBAL.geodeRobot = max(GLOBAL.geodeRobot, buildGeoRobot.geodeRobot)

		ret = max(ret, buildGeoRobot.nextGeodes(b, maxSteps))
	}

	if canBuildObs && s.obsidianRobot < maxObsPerMin && s.obsidian < 2*b.GeodeRobotCostsObsidian {
		// &&
		//(s.obsidianRobot == 0 || b.
		buildObsRobot := s
		buildObsRobot.ore -= b.ObsidianRobotCostsOre
		buildObsRobot.clay -= b.ObsidianRobotCostsClay
		buildObsRobot.MinutePassesCollect()
		buildObsRobot.obsidianRobot++
		GLOBAL.obsidianRobot = max(GLOBAL.obsidianRobot, buildObsRobot.obsidianRobot)
		ret = max(ret, buildObsRobot.nextGeodes(b, maxSteps))
	}

	if !canBuildGeo {
		if canBuildOre && s.oreRobot < maxOrePerMin && s.ore < 2*maxOrePerMin {
			buildOreRobot := s
			buildOreRobot.ore -= b.OreRobotCost
			buildOreRobot.MinutePassesCollect()
			buildOreRobot.oreRobot++
			GLOBAL.oreRobot = max(GLOBAL.oreRobot, buildOreRobot.oreRobot)
			ret = max(ret, buildOreRobot.nextGeodes(b, maxSteps))
		}

		if canBuildClay && s.clayRobot < maxClayPerMin && s.clay < 2*maxClayPerMin {
			buildClayRobot := s
			buildClayRobot.ore -= b.ClayRobotCost
			buildClayRobot.MinutePassesCollect()
			buildClayRobot.clayRobot++
			GLOBAL.clayRobot = max(GLOBAL.clayRobot, buildClayRobot.clayRobot)
			ret = max(ret, buildClayRobot.nextGeodes(b, maxSteps))
		}

		// if !(canBuildOre && canBuildClay)
		if s.ore < b.OreRobotCost || s.ore < b.ClayRobotCost || s.ore < b.ObsidianRobotCostsOre || s.ore < b.GeodeRobotCostsOre ||
			s.clay < b.ObsidianRobotCostsClay {
			cantBuildAny := s
			cantBuildAny.MinutePassesCollect()
			// GLOBAL.oreRobot = max(GLOBAL.oreRobot, buildOreRobot.oreRobot)
			ret = max(ret, cantBuildAny.nextGeodes(b, maxSteps))
		}
	}

	return ret
}
