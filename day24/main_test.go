package main

import (
	"container/heap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeap(t *testing.T) {

	q := &IntHeap{dest: position{5, 5}}

	heap.Push(q, move{to: position{1, 0}, stepsTo: 1})

	heap.Push(q, move{to: position{3, 0}, stepsTo: 15})
	heap.Push(q, move{to: position{3, 0}, stepsTo: 3})

	v1 := heap.Pop(q).(move)

	assert.Equal(t, position{3, 0}, v1.to)
	assert.Equal(t, 3, v1.stepsTo)

	v2 := heap.Pop(q).(move)

	assert.Equal(t, position{3, 0}, v2.to)
	assert.Equal(t, 15, v2.stepsTo)

	v3 := heap.Pop(q).(move)

	assert.Equal(t, position{1, 0}, v3.to)
	assert.Equal(t, 1, v3.stepsTo)
}

func TestRunBlizzards(t *testing.T) {
	xLen, yLen := 5, 5
	sim := simulation{
		blizzardStart{position{0, 0}, Right, &xLen, &yLen},
		blizzardStart{position{4, 0}, Left, &xLen, &yLen},
		blizzardStart{position{0, 0}, Down, &xLen, &yLen},
		blizzardStart{position{0, 4}, Up, &xLen, &yLen},
	}

	t4 := sim.RunBlizzards(4)
	assert.Equal(t, position{4, 0}, t4[0].start)
	assert.Equal(t, position{0, 0}, t4[1].start)
	assert.Equal(t, position{0, 4}, t4[2].start)
	assert.Equal(t, position{0, 0}, t4[3].start)

	assert.Equal(t, true, t4.IsStorming(position{0, 0}))
	assert.Equal(t, false, t4.IsStorming(position{1, 1}))

	//Wrap Around
	tWrapped := sim.RunBlizzards(5)
	assert.Equal(t, position{0, 0}, tWrapped[0].start)
	assert.Equal(t, position{4, 0}, tWrapped[1].start)
	assert.Equal(t, position{0, 0}, tWrapped[2].start)
	assert.Equal(t, position{0, 4}, tWrapped[3].start)

	assert.Equal(t, true, tWrapped.IsStorming(position{0, 0}))
	assert.Equal(t, false, tWrapped.IsStorming(position{1, 1}))
}
