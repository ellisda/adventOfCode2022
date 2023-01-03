package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"cd.splunkdev.com/genuine-cat/advent-of-code-2022/davide/day20/pkg/linkedlist"
	"github.com/samber/lo"
)

type circ interface {
	Ints() []int
	// Copy() interface{}
	CycleN(int)
	MoveItem(int)
	Len() int
	FindItemsFromZero(indices ...int) []int
	Distr() map[int]int
}

func main() {
	in := "input.txt"
	file, _ := os.Open(in)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	nums := make([]int, 0)
	for scanner.Scan() {
		n, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal("parse fail")
		}
		nums = append(nums, n)
	}
	c := NewCircular(nums)
	fmt.Println("Read len", len(c.Ints()))

	// fmt.Println("Before dist - len", c.Len(), "dist", c.Distr())
	c.CycleN(1)

	// fmt.Println("After dist - len", c.Len(), "dist", c.Distr())

	vals := c.FindItemsFromZero(1000, 2000, 3000)
	fmt.Println("Part1", vals, "sum", lo.Sum(vals))

	decrypted := lo.Map(nums, func(v int, _ int) int {
		ret := v * 811589153
		if (ret > 0) != (v > 0) {
			panic("overflow")
		}
		return ret
	})
	c2 := NewCircular(decrypted)

	c2.CycleN(10)

	vals2 := c2.FindItemsFromZero(1000, 2000, 3000)
	fmt.Println("Part2", vals2, "sum", lo.Sum(vals2))

}

func NewCircular(s []int) circ {
	// ret := myslice.NewCircular(s)
	// return &ret

	return linkedlist.NewCircular(s)
}
