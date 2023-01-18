package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJetSeq(t *testing.T) {
	j2 := parseJetSequence(INPUT2)

	assert.Equal(t, false, j2.NextIsLeft())
	assert.Equal(t, false, j2.NextIsLeft())
	assert.Equal(t, false, j2.NextIsLeft())
	assert.Equal(t, true, j2.NextIsLeft())
	assert.Equal(t, true, j2.NextIsLeft())
	assert.Equal(t, false, j2.NextIsLeft())
	assert.Equal(t, true, j2.NextIsLeft())
}
