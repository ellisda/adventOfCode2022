package myslice

import (
	"fmt"
	"log"
)

func NewCircular(s []int) circular {
	c := make(circular, len(s))
	for i := 0; i < len(s); i++ {
		c[i] = &item{value: s[i], origIndex: i}
	}
	return c
}

type circular []*item

type item struct {
	value     int
	origIndex int
}

func (c circular) Copy() *circular {
	ret := make(circular, len(c))
	if n := copy(ret, c); n != len(c) {
		log.Fatal("copy failure")
	}

	return &ret
}

func (c circular) Len() int {
	return len(c)
}

func (c circular) Ints() []int {
	ret := make([]int, len(c))
	for i, v := range c {
		ret[i] = v.value
	}
	return ret
}

// CycleOnce returns a decrypted sequence
func (c circular) CycleOnce() {
	orig := c.Copy()
	ret := c.Copy()
	//We only have to scan 5000 items and each only moves at most ~5000 spaces (w/ wrapping)
	for _, selected := range *orig {
		i := c.FindIndexOf(selected)
		if i < 0 {
			log.Fatal("bad index")
		}

		MoveItemFunc(ret, i)
	}
	// return ret
}

func (c circular) FindItemsFromZero(indices ...int) []int {
	return nil
}

func (c circular) FindIndexOf(match *item) int {
	for i := range c {
		if c[i] == match {
			return i
		}
	}
	return -1
}

func (c *circular) MoveItem(i int) {
	MoveItemFunc(c, i)
}

func MoveItemFunc(ret *circular, i int) {
	v := (*ret)[i]
	next := (i + v.value) % len(*ret)

	if v.value >= 0 {
		// left0 := (*ret)[:v] //still needs
		//Move the inner items left to allow v to move
		copy((*ret)[i:next], (*ret)[i+1:next+1])

		(*ret)[next] = v
	} else {
		next += len(*ret)
		fmt.Println("dest", (*ret)[next:i])
		fmt.Println("src", (*ret)[next:i+1])
		copy((*ret)[next:i], (*ret)[next:i+1])

		(*ret)[next] = v

	}

	// right := append((*ret)[:i]NewCircular([]int{(*ret)[i]})// dd :=ret[i]
	// (*ret) = append((*ret)[:i], (*ret)[i:]...)
	// _ = i

}
