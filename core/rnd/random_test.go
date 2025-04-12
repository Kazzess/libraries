package rnd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	str := RandomString(100)
	assert.Equal(t, 100, len(str))
}

func BenchmarkRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandomString(100)
	}
}

func TestRandomCase(t *testing.T) {
	cases := []interface{}{"1", "abc", "qwerty", "2", "3", "0"}
	res := RandomCase(cases...)
	assert.True(t, func() bool {
		for _, v := range cases {
			if v == res {
				return true
			}
		}

		return false
	}())
}

func BenchmarkRandomCase(b *testing.B) {
	cases := []interface{}{"1", "abc", "qwerty", "2", "3", "0"}

	for i := 0; i < b.N; i++ {
		RandomCase(cases...)
	}
}
