package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RandomStringNotEmpty(t *testing.T) {
	last := ""
	for c := 0; c < 1000; c++ {
		s := randString(10)
		assert.True(t, s != "")
		assert.True(t, len(s) == 10)
		assert.NotEqual(t, last, s)
		last = s
	}
}
