package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRotate(t *testing.T) {
	assert.Equal(t, DOWN, RIGHT.Rotate(true))
	assert.Equal(t, LEFT, DOWN.Rotate(true))
	assert.Equal(t, UP, LEFT.Rotate(true))
	assert.Equal(t, RIGHT, UP.Rotate(true))

	assert.Equal(t, UP, RIGHT.Rotate(false))
	assert.Equal(t, RIGHT, DOWN.Rotate(false))
	assert.Equal(t, DOWN, LEFT.Rotate(false))
	assert.Equal(t, LEFT, UP.Rotate(false))
}
