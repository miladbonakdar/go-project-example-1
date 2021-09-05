package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestMin(t *testing.T) {
	t.Parallel()
	cases := []struct {
		a, b, result int
	}{
		{a: 1, b: 2, result: 1},
		{a: 10, b: 2, result: 2},
	}
	for _, item := range cases {
		assert.Equal(t, item.result, Min(item.a, item.b))
	}
}

func TestMinUint(t *testing.T) {
	t.Parallel()
	cases := []struct {
		a, b, result uint
	}{
		{a: 1, b: 2, result: 1},
		{a: 10, b: 2, result: 2},
	}
	for _, item := range cases {
		assert.Equal(t, item.result, MinUint(item.a, item.b))
	}
}

func TestMaxUint(t *testing.T) {
	t.Parallel()
	cases := []struct {
		a, b, result uint
	}{
		{a: 1, b: 2, result: 2},
		{a: 10, b: 2, result: 10},
	}
	for _, item := range cases {
		assert.Equal(t, item.result, MaxUint(item.a, item.b))
	}
}
