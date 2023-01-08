package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoveN(t *testing.T) {

	g := gang{
		elf{pos{1, 2}},
		elf{pos{2, 2}},
	}

	g1 := g.MoveN(1)

	assert.Equal(t, pos{1, 1}, g1[0].loc)
	assert.Equal(t, pos{2, 1}, g1[1].loc)

	g2 := g.MoveN(2)

	assert.Equal(t, g[0].loc, g2[0].loc)
	assert.Equal(t, g[1].loc, g2[1].loc)
}
