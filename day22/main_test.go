package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	assert.Equal(t, DOWN, RIGHT.Rotate(CLOCKWISE))
	assert.Equal(t, LEFT, DOWN.Rotate(CLOCKWISE))
	assert.Equal(t, UP, LEFT.Rotate(CLOCKWISE))
	assert.Equal(t, RIGHT, UP.Rotate(CLOCKWISE))

	assert.Equal(t, UP, RIGHT.Rotate(COUNTERCLOCKWISE))
	assert.Equal(t, RIGHT, DOWN.Rotate(COUNTERCLOCKWISE))
	assert.Equal(t, DOWN, LEFT.Rotate(COUNTERCLOCKWISE))
	assert.Equal(t, LEFT, UP.Rotate(COUNTERCLOCKWISE))

	assert.Equal(t, UP, UP.Rotate(NOROTATION))
	assert.Equal(t, DOWN, DOWN.Rotate(NOROTATION))
	assert.Equal(t, RIGHT, RIGHT.Rotate(NOROTATION))
	assert.Equal(t, LEFT, LEFT.Rotate(NOROTATION))
}

func TestR1MinXLeft(t *testing.T) {
	g, _ := parseInput(strings.Split(INPUT, "\n"), true)

	assert.Equal(t, player{pos{0, 149}, RIGHT}, g.rows[0].minXNext)

	assert.Equal(t, player{pos{0, 100}, RIGHT}, g.rows[49].minXNext)
}
