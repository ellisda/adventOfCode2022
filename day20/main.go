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
	CycleOnce()
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
	c.CycleOnce()

	// fmt.Println("After dist - len", c.Len(), "dist", c.Distr())

	vals := c.FindItemsFromZero(1000, 2000, 3000)
	fmt.Println("Part1", vals, "sum", lo.Sum(vals))

}

func NewCircular(s []int) circ {
	// ret := myslice.NewCircular(s)
	// return &ret

	return linkedlist.NewCircular(s)
}
