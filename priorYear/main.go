// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"io"
	"math"
	"os"
)

func main() {

	f, err := os.Open("input.txt")
	if err != nil {

	}

	// problem1(f)
	problem2(f)
	f.Close()
	f, err = os.Open("input.txt")
	problem2Stack(f)

}

func problem1(f io.Reader) {
	// seed the `prev` with a large number to ensure the first line isn't counted as an increase
	increases, prev := 0, math.MaxInt64
	for {
		var i int
		n, err := fmt.Fscanln(f, &i)

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}

		if i > prev {
			increases++
		}
		fmt.Println(i, n, increases)

		prev = i

	}
	fmt.Println("Increases: ", increases)

}

func problem2(f io.Reader) {
	// seed the `prev` with a large number to ensure the first line isn't counted as an increase
	increases, prevSum3 := 0, 0.0
	linesRead := 0

	for {
		var i int
		n, err := fmt.Fscanln(f, &i)

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}
		linesRead++

		if linesRead < 3 {
			prevSum3 += float64(i)
		} else if linesRead == 3 {
			prevSum3 = (prevSum3 + float64(i)) / 3
		} else { //only compare when we have consecutive groups of 3
			sum3 := newSum3(prevSum3, i)
			if sum3 > prevSum3 {
				increases++
			}
			prevSum3 = sum3
		}
		fmt.Println(i, n, "[", linesRead, linesRead%3, "]", prevSum3, prevSum3*3, increases)
	}
	fmt.Println("Increases: ", increases)

}

func newSum3(oldSum float64, newValue int) float64 {
	return (oldSum*3 + float64(newValue)) / 4
}

type stack []int

func (s stack) Sum() int {
	ret := 0
	for _, si := range s {
		ret += si
	}
	return ret
}

func problem2Stack(f io.Reader) {

	s := make(stack, 3)
	// seed the `prev` with a large number to ensure the first line isn't counted as an increase
	increases, prevSum3 := 0, 0
	linesRead := 0

	for {
		var i int
		n, err := fmt.Fscanln(f, &i)

		// end of the line (EOF) or not
		if n == 0 || err == io.EOF {
			break
		}
		linesRead++

		if linesRead < 3 {
			s[linesRead-1] = i
		} else if linesRead == 3 {
			s[linesRead-1] = i
			prevSum3 = s.Sum()
		} else { //only compare when we have consecutive groups of 3
			s = append(s[1:], i)
			sum3 := s.Sum()
			if sum3 > prevSum3 {
				increases++
			}
			prevSum3 = sum3
		}
		// fmt.Println(i, n, "[", linesRead, linesRead%3, "]", s, prevSum3, increases)
	}
	fmt.Println("Increases: ", increases)

}
