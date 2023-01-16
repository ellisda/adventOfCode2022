package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZMin(t *testing.T) {
	s2 := parseInput(strings.Split(INPUT2, "\n"))

	water := FillWithWater(s2)

	//in z == 1 plane, around point 2,2,1
	assert.False(t, s2.IsEmpty(pos{2, 2, 1}))

	assert.False(t, water[pos{2, 2, 1}])

	assert.True(t, water[pos{2, 2, 0}])

	//Make sure Neighbors in the same z=1 plane are all water
	assert.True(t, water[pos{2, 1, 1}])
	assert.True(t, water[pos{1, 2, 1}])

	assert.True(t, water[pos{2, 3, 1}])
	assert.True(t, water[pos{3, 2, 1}])

	assert.True(t, water[pos{1, 1, 1}])
	assert.True(t, water[pos{0, 1, 1}])
	assert.True(t, water[pos{1, 0, 1}])
}

func TestZMax(t *testing.T) {
	s2 := parseInput(strings.Split(INPUT2, "\n"))

	water := FillWithWater(s2)

	p := pos{2, 2, 6}
	//in z == 1 plane, around point 2,2,1
	assert.False(t, s2.IsEmpty(p))

	assert.False(t, water[p])

	assert.True(t, water[pos{p.x, p.y, p.z + 1}])

	//Make sure Neighbors in the same z=1 plane are all water
	assert.True(t, water[pos{p.x - 1, p.y, p.z}])
	assert.True(t, water[pos{p.x, p.y - 1, p.z}])

	assert.True(t, water[pos{p.x + 1, p.y, p.z}])
	assert.True(t, water[pos{p.x, p.y + 1, p.z}])
}
