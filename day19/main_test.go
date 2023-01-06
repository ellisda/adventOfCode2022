package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStep(t *testing.T) {
	bp1 := blueprint{1, 4, 2, 3, 14, 2, 7}
	bp2 := blueprint{2, 2, 3, 3, 8, 3, 12}

	s1 := state{oreRobot: 1, ore: 0, step: 0}

	//After min 1
	g2 := s1.nextGeodes(&bp1, 2)
	assert.Equal(t, 0, g2)

	assert.Equal(t, 9, s1.nextGeodes(&bp1, 24))

	assert.Equal(t, 12, s1.nextGeodes(&bp2, 24))
}

func TestCanMakeOne(t *testing.T) {
	bp := blueprint{6, 4, 4, 3, 7, 4, 11}
	s1 := state{oreRobot: 1, ore: 0, step: 0}

	assert.Equal(t, 4, s1.nextGeodes(&bp, 24))
}

func TestPart2_partial(t *testing.T) {
	bp1 := blueprint{1, 4, 2, 3, 14, 2, 7}

	s1 := state{oreRobot: 1, ore: 0, step: 0}

	// assert.Equal(t, 15, s1.nextGeodes(&bp1, 26))
	// assert.Equal(t, 20, s1.nextGeodes_part2(&bp1, 27))
	// assert.Equal(t, 26, s1.nextGeodes_part2(&bp1, 28))
	// assert.Equal(t, 32, s1.nextGeodes(&bp1, 29))
	assert.Equal(t, 39, s1.nextGeodes(&bp1, 30))
	// assert.Equal(t, 47, s1.nextGeodes_part2(&bp1, 31))
	assert.Equal(t, 56, s1.nextGeodes(&bp1, 32))

	assert.Equal(t, state{}, GLOBAL)
}

func TestPart2_bp2(t *testing.T) {
	bp2 := blueprint{2, 2, 3, 3, 8, 3, 12}

	s1 := state{oreRobot: 1, ore: 0, step: 0}

	assert.Equal(t, 62, s1.nextGeodes_part2(&bp2, 32))

}
