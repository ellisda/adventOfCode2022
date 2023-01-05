package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStep(t *testing.T) {
	bp := blueprint{1, 4, 2, 3, 14, 2, 7}

	//After min 1
	s1 := state{oreRobot: 1, ore: 0, step: 0}
	g2 := s1.nextGeodes(&bp, 2)
	assert.Equal(t, 0, g2)

	assert.Equal(t, 1, s1.nextGeodes(&bp, 24))

}
