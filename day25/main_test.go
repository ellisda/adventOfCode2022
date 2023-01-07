package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CarryOnes(t *testing.T) {
	s := snafu{3}
	s.CarryOnes()

	assert.Equal(t, 1, s[1])
	assert.Equal(t, 2, s[0])
}

func Test_CarryOnesBorrow(t *testing.T) {
	s := snafu{-3}
	s.CarryOnes()

	assert.Equal(t, -1, s[1])
	assert.Equal(t, 2, s[0])
}

func Test_CarryOnesExample(t *testing.T) {
	s := snafu{10, 11, -2, 4, 2, 1}
	s.CarryOnes()

	exp := snafu{0, -2, 1, -1, -2, 2}
	assert.Equal(t, exp, s)
	assert.Equal(t, "2=-1=0", exp.String())
}

func Test_Add(t *testing.T) {
	s := snafu{1, 1, 1, 1, 1, 1}
	s2 := snafu{0, 1, 1, 1, 1, 1}

	exp := snafu{1, 2, 2, 2, 2, 2}
	assert.Equal(t, exp, s.Add(s2))
}

func Test_AddNegative(t *testing.T) {
	s := snafu{1, 1, 1, 1, 1, 1}
	s2 := snafu{0, 1, -1, 1, -2, 1}

	exp := snafu{1, 2, 0, 2, -1, 2}
	assert.Equal(t, exp, s.Add(s2))
}
