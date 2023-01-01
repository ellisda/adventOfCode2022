package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"

	"github.com/samber/lo"
)

type instruction struct {
	sensor position
	beacon position
}

type position struct {
	x int
	y int
}

func main() {
	in := "input.txt"
	file, _ := os.Open(in)
	defer file.Close()

	input := []instruction{}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		iv := scanner.Text()
		ins := instruction{}
		f := `Sensor at x=%d, y=%d: closest beacon is at x=%d, y=%d`
		n, err := fmt.Sscanf(iv, f, &ins.sensor.x, &ins.sensor.y, &ins.beacon.x, &ins.beacon.y)
		if n < 4 || err != nil {
			log.Fatal(err)
		}

		input = append(input, ins)
		fmt.Println(ins)
	}

	if in == "input.txt" {
		fmt.Println("Part1 (input.txt)", countEmpties(2000000, input))
		fmt.Print("Part2 (input.txt)")
		countEmptiesWithinRange(input, 0, 4000000)

	} else {
		fmt.Println("Part1 (input2.txt)", countEmpties(10, input))
		fmt.Print("Part2 (input2.txt)")
		countEmptiesWithinRange(input, 0, 20)
	}
}

func countEmpties(row int, instrs []instruction) int {
	m := make(map[int]bool)
	for _, ins := range instrs {
		//Step 1 - determine is sensor is close enough to row to have any x moves
		//Step 2 - store map of x positions where sensor "saw" no beacons
		//Final  - return count of seen positions from map
		//    --- make sure to check that no beacon exists at each of the "seen" x positions

		y_dist := abs(row - ins.sensor.y)
		x_Remain := ins.distance() - y_dist
		fmt.Println("Sensor", ins.sensor, "y_dist", y_dist, "sensor_dist", ins.distance(), "x_remain", x_Remain)
		for x := 0; x <= x_Remain; x++ {
			m[ins.sensor.x+x] = true
			m[ins.sensor.x-x] = true
		}

	}
	for _, ins := range instrs {
		if ins.beacon.y == row {
			delete(m, ins.beacon.x)
		}
	}
	seq := lo.Keys(m)
	sort.Ints(seq)
	// fmt.Println("Seen wo/beacons", seq)

	return len(m)
}

func abs(a int) int {
	return int(math.Abs(float64(a)))
}

// distance returns the manhatten distance from the closest beacon
func (i instruction) distance() int {
	return abs(i.beacon.x-i.sensor.x) + abs(i.beacon.y-i.sensor.y)
}

func countEmptiesWithinRange(instrs []instruction, min, max int) {
	for y := min; y <= max; y++ {
		// fmt.Println("y", y)
	GRID:
		for x := min; x <= max; x++ {
			for _, ins := range instrs {

				//if position is closer than closest beacon, we can skip ahead till end of current sensor range
				if d := ins.distanceFromSensor(x, y); d <= ins.distance() {
					nextX := ins.getNextUnseenX(x, y)
					// fmt.Println("ins", ins, "is", d, "from", x, y, "which is closer than it's beacon's dist", ins.distance())
					// fmt.Println("ins", ins, "sees position", x, y, "next unseen x", nextX)

					if nextX > 0 {
						x = nextX - 1 //-1 to allow for-loop increment
					}
					continue GRID
				}

			}
			fmt.Println("Found Empty at", x, y, "freq", 4000000*x+y)
		}
	}
}

func (i instruction) distanceFromSensor(targetX, targetY int) int {
	return abs(targetX-i.sensor.x) + abs(targetY-i.sensor.y)
}

func (ins instruction) getNextUnseenX(x, y int) int {
	y_dist := abs(y - ins.sensor.y)
	x_Remain := ins.distance() - y_dist
	if x_Remain <= 0 {
		return -1
	}
	xLeftVis := ins.sensor.x - x_Remain
	xRightVis := ins.sensor.x + x_Remain
	if x >= xLeftVis && x <= xRightVis {
		return xRightVis + 1
	}
	return -1
}
