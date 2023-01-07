package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type snafu [30]int

func main() {

	in := "input.txt"
	file, _ := os.Open(in)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	fmt.Println(addLines(lines...))

}

func addLines(lines ...string) string {
	snafu_sum := snafu{}

	for _, l := range lines {
		s := parseSNAFU(l)
		snafu_sum = snafu_sum.Add(s)
		fmt.Println()
	}

	snafu_sum.CarryOnes()
	return snafu_sum.String()
}

func abs(a int) int {
	if a > 0 {
		return a
	}
	return -a
}

func parseSNAFU(l string) snafu {
	ret := snafu{}
	sPos := 0 // least significant bit goes in s[0]
	for pos := len(l) - 1; pos >= 0; pos-- {
		switch l[pos] {
		case '1':
			ret[sPos]++
		case '2':
			ret[sPos] += 2
		case '0':
		case '-':
			ret[sPos]--

		case '=':
			ret[sPos] -= 2
		}
		sPos++
	}
	return ret
}

func (s snafu) Add(b snafu) snafu {
	for pos := 0; pos < len(s); pos++ {
		s[pos] += b[pos]
	}
	return s
}

func (s *snafu) CarryOnes() {
	for pos := 0; pos < len(*s); pos++ {
		if carry := s[pos] / 5; carry != 0 {
			s[pos+1] += carry
			s[pos] = s[pos] % 5
		}
		if v := s[pos]; v > 2 {
			s[pos+1]++
			s[pos] = v - 5
		} else if v < -2 {
			s[pos+1]--
			s[pos] = v + 5
		}
	}
}

func (s *snafu) String() string {
	b := strings.Builder{}

	msb := -1
	for pos := len(s) - 1; pos >= 0; pos-- {
		if s[pos] != 0 {
			msb = pos
			break
		}
	}

	for pos := msb; pos >= 0; pos-- {
		switch s[pos] {
		case -2:
			b.WriteRune('=')
		case -1:
			b.WriteRune('-')
		case 0:
			b.WriteRune('0')
		case 1:
			b.WriteRune('1')
		case 2:
			b.WriteRune('2')
		}
	}

	return b.String()
}
