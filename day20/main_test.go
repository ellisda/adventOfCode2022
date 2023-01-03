package main

import (
	"testing"

	"cd.splunkdev.com/genuine-cat/advent-of-code-2022/davide/day20/pkg/linkedlist"
	"github.com/stretchr/testify/assert"
)

var testInput = []int{
	1, 2, -3, 3, -2, 0, 4,
}

func TestMoveCircular(t *testing.T) {
	c := NewCircular(testInput)
	// orig := c.Copy()

	c.MoveItem(1)

	assert.Equal(t, []int{
		1, -3, 3, 2, -2, 0, 4,
	}, c.Ints())
}

func TestMoveCircularBackwards(t *testing.T) {
	c := NewCircular(testInput)

	c.MoveItem(2)

	assert.Equal(t, []int{
		1, 2, 3, -2, 0, -3, 4,
	}, c.Ints())
}

func TestMoveCircularSingleForward(t *testing.T) {
	c := linkedlist.NewCircular(testInput)

	c.MoveItem(0)

	rotated := c.FindItem(c.Len() - 1)
	assert.Equal(t, []int{
		2, 1, -3, 3, -2, 0, 4,
	}, rotated.Ints())
}

func TestMoveCycleOnce_Steps(t *testing.T) {
	c := linkedlist.NewCircular(testInput)

	s1, s2, s3, s4, s5, s6, s7 := c, c.FindItem(1), c.FindItem(2), c.FindItem(3), c.FindItem(4), c.FindItem(5), c.FindItem(6)
	s1.MoveItem(0)

	rotated1 := c.FindItem(c.Len() - 1)
	assert.Equal(t, []int{
		2, 1, -3, 3, -2, 0, 4,
	}, rotated1.Ints())

	s2.MoveItem(0)
	rotated2 := s2.FindItem(c.Len() - 2)
	assert.Equal(t, []int{
		1, -3, 2, 3, -2, 0, 4,
	}, rotated2.Ints())

	s3.MoveItem(0)
	rotated3 := s3.FindItem(c.Len() - 4)
	assert.Equal(t, []int{
		1, 2, 3, -2, -3, 0, 4,
	}, rotated3.Ints())

	s4.MoveItem(0)
	rotated4 := s4.FindItem(c.Len() + 2)
	assert.Equal(t, []int{
		1, 2, -2, -3, 0, 3, 4,
	}, rotated4.Ints())

	s5.MoveItem(0)
	rotated5 := s5.FindItem(c.Len() + 1)
	assert.Equal(t, []int{
		1, 2, -3, 0, 3, 4, -2,
	}, rotated5.Ints())

	s6.MoveItem(0)
	rotated6 := s6.FindItem(c.Len() - 3)
	assert.Equal(t, []int{
		1, 2, -3, 0, 3, 4, -2,
	}, rotated6.Ints())

	s7.MoveItem(0)
	rotated7 := s7.FindItem(c.Len() + 4)
	assert.Equal(t, []int{
		1, 2, -3, 4, 0, 3, -2,
	}, rotated7.Ints())

}

func TestCycleOnce(t *testing.T) {
	c := NewCircular(testInput)

	c.CycleOnce()

	assert.Equal(t, []int{
		1, 2, -3, 4, 0, 3, -2,
	}, c.Ints())
}

func TestBackOne(t *testing.T) {
	c := NewCircular([]int{
		1, 2, -3, 3, -1, 0, 4,
	})

	c.MoveItem(4)

	assert.Equal(t, []int{
		1, 2, -3, -1, 3, 0, 4,
	}, c.Ints())
}

func TestMoveAllTheWayAround_DoNothing(t *testing.T) {
	input := []int{
		5, -5, 5, -5, 5, -5,
	}
	c := NewCircular(input)

	c.MoveItem(0)
	assert.Equal(t, input, c.Ints())

	c.MoveItem(1)
	assert.Equal(t, input, c.Ints())

	c.CycleOnce()
	assert.Equal(t, input, c.Ints())
}

func TestLargeNum(t *testing.T) {
	// 5292 % 5 == 2
	c := linkedlist.NewCircular([]int{
		5292, 4025, -5805, -6388, -6433, -2584,
	})

	c.MoveItem(0)

	assert.Equal(t, []int{
		4025, -5805, 5292, -6388, -6433, -2584,
	}, c.FindItem(-2).Ints())

}
