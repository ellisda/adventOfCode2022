package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/samber/lo"
)

const (
	Mult             operation               = '*'
	Add              operation               = '+'
	SelfTestArgument testOperationArgumenent = -1
)

type operation rune
type testOperationArgumenent int

type monkey struct {
	items       []int
	op          operation
	opArg       testOperationArgumenent
	testDivisor int
	testPassIdx int
	testFaiIdx  int
	inspections int
}

type monkeys []*monkey

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		panic(err)
	}

	monkeys := parseInput(f)
	fmt.Println(monkeys)

	for i := 0; i < 20; i++ {
		for _, m := range monkeys {
			for range m.items {
				m.ThrowTo(monkeys, true, 1)
			}
		}
	}
	fmt.Println()
	fmt.Println("Part 1 -- Monkeys are at")
	fmt.Println(monkeys.String())
	counts := lo.Map(monkeys, func(m *monkey, index int) int { return m.inspections })
	sort.Ints(counts)
	fmt.Println(counts)
	fmt.Println("goof", counts[len(counts)-1]*counts[len(counts)-2])

	f.Seek(0, io.SeekStart)
	monkeys = parseInput(f)
	fmt.Println(monkeys)
	commonDiv := 1
	for _, m := range monkeys {
		commonDiv = commonDiv * m.testDivisor
	}

	for i := 0; i < 10000; i++ {
		for _, m := range monkeys {
			for range m.items {
				m.ThrowTo(monkeys, false, commonDiv)
			}
		}
	}

	fmt.Println()
	fmt.Println("Part 2 -- Monkeys are at")
	fmt.Println(monkeys.String())

	fmt.Println()
	fmt.Println(monkeys.String())
	counts = lo.Map(monkeys, func(m *monkey, index int) int { return m.inspections })
	sort.Ints(counts)
	fmt.Println(counts)
	fmt.Println("goof", counts[len(counts)-1]*counts[len(counts)-2])
}

func (ms monkeys) String() string {
	s := ""
	for _, m := range ms {
		s += fmt.Sprintf("%v, %d; ", m.items, m.inspections)
	}
	return s
}

func (m *monkey) ThrowTo(monkeys monkeys, part1 bool, commonDiv int) {
	if len(m.items) == 0 {
		return
	}

	// fmt.Println("Before: m.items", m.items)
	item := m.items[0]

	var dest *monkey
	item = m.InspectItem(item)
	if part1 {
		item = item / 3
	} else {
		item = item % commonDiv
	}
	if item%m.testDivisor == 0 {
		dest = monkeys[m.testPassIdx]
	} else {
		dest = monkeys[m.testFaiIdx]
	}

	dest.items = append(dest.items, item)
	m.items = m.items[1:]
	// fmt.Println("After: m.items", m.items, "dest: ", dest.items)
}

func (m *monkey) InspectItem(item int) (worry int) {
	m.inspections++
	v := int(m.opArg)
	if v == int(SelfTestArgument) {
		v = item
	}
	switch m.op {
	case Mult:
		return item * v
	case Add:
		return item + v
	}
	panic("no case?!")
}

func parseInput(f io.ReadSeekCloser) monkeys {
	ret := monkeys{}
	i := 0
	s := bufio.NewScanner(f)
	s.Split(bufio.ScanLines)
	m := &monkey{}
	for s.Scan() {
		fmt.Println(s.Text())
		switch i % 7 {
		case 0:
			//monkey
			m = &monkey{}
		case 1:
			index := len("  Starting items: ")
			items := strings.Split(s.Text()[index:], ", ")

			m.items = lo.Map(items, func(s string, n int) int {
				i, _ := strconv.Atoi(s)
				return i
			})
		case 2:
			if s.Text() == "  Operation: new = old * old" {
				m.op = operation(Mult)
				m.opArg = SelfTestArgument
			} else if n, err := fmt.Sscanf(s.Text(), "  Operation: new = old %c %d", &m.op, &m.opArg); err != nil || n < 2 {
				log.Fatal("parse fail")
			}
		case 3:
			if n, err := fmt.Sscanf(s.Text(), "  Test: divisible by %d", &m.testDivisor); err != nil || n < 1 {
				log.Fatal("parse fail")
			}
		case 4:
			if n, err := fmt.Sscanf(s.Text(), "    If true: throw to monkey %d", &m.testPassIdx); err != nil || n < 1 {
				log.Fatal("parse fail")
			}
		case 5:
			if n, err := fmt.Sscanf(s.Text(), "    If false: throw to monkey %d", &m.testFaiIdx); err != nil || n < 1 {
				log.Fatal("parse fail")
			}

		case 6:
			ret = append(ret, m)
		}

		i++
	}

	ret = append(ret, m)
	return ret
}
