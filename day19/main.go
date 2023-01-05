package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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
	in := "input2.txt"
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

	for _, bp := range bps {
		fmt.Println("blueprint", bp, "best", bp.MaxGeodes(23))
	}
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

// param is a pointer cause this is a hot path and less memory alloc
func (s state) nextGeodes(b *blueprint, maxSteps int) int {
	if s.step >= maxSteps {
		return s.geode
	}

	s.ore += s.oreRobot
	s.clay += s.clayRobot
	s.obsidian += s.obsidianRobot
	s.geode += s.geodeRobot
	s.step++
	ret := s.geode

	if buildGeoRobot := s; buildGeoRobot.ore >= b.GeodeRobotCostsOre && buildGeoRobot.obsidian >= b.GeodeRobotCostsObsidian {
		buildGeoRobot.ore -= b.GeodeRobotCostsOre
		buildGeoRobot.obsidian -= b.GeodeRobotCostsObsidian
		buildGeoRobot.geodeRobot++

		fmt.Println("Build Geo Robot", buildGeoRobot)

		ret = max(ret, buildGeoRobot.nextGeodes(b, maxSteps))
	} else {
		if buildObsRobot := s; buildObsRobot.ore >= b.ObsidianRobotCostsOre && buildObsRobot.clay >= b.ObsidianRobotCostsClay {
			buildObsRobot.ore -= b.ObsidianRobotCostsOre
			buildObsRobot.clay -= b.ObsidianRobotCostsClay
			buildObsRobot.obsidianRobot++
			ret = max(ret, buildObsRobot.nextGeodes(b, maxSteps))
		} else {
			if buildOreRobot := s; buildOreRobot.ore >= b.OreRobotCost {
				buildOreRobot.ore -= b.OreRobotCost
				buildOreRobot.oreRobot++
				ret = max(ret, buildOreRobot.nextGeodes(b, maxSteps))
			}

			if buildClayRobot := s; buildClayRobot.ore >= b.ClayRobotCost {
				buildClayRobot.ore -= b.ClayRobotCost
				buildClayRobot.clayRobot++
				ret = max(ret, buildClayRobot.nextGeodes(b, maxSteps))
			}

			if cantBuildAny := s; cantBuildAny.ore < b.OreRobotCost &&
				cantBuildAny.ore < b.ObsidianRobotCostsOre &&
				cantBuildAny.ore < b.GeodeRobotCostsOre &&
				cantBuildAny.ore < b.ClayRobotCost {
				ret = max(ret, cantBuildAny.nextGeodes(b, maxSteps))
			}
		}

	}

	return ret
}
