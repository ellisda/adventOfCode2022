package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestR1MinXLeft(t *testing.T) {
	s2 := parseInput(strings.Split(INPUT2, "\n"))

	water := FillWithWater(s2)

	//in z == 1 plane, around point 2,2,1
	assert.False(t, s2.IsEmpty(pos{2, 2, 1}))

	assert.False(t, water[pos{2, 2, 1}])

	assert.True(t, water[pos{2, 2, 0}])

	assert.True(t, water[pos{2, 1, 1}])
	assert.True(t, water[pos{1, 2, 1}])

	assert.True(t, water[pos{2, 3, 1}])
	assert.True(t, water[pos{3, 2, 1}])

	assert.True(t, water[pos{1, 1, 1}])
	assert.True(t, water[pos{0, 1, 1}])
	assert.True(t, water[pos{1, 0, 1}])
}
