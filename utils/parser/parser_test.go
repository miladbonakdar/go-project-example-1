package parser_test

import (
	"giftcard-engine/utils/parser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseNumber(t *testing.T) {
	t.Parallel()
	cases := []struct {
		stringNumber string
		number       uint
	}{
		{stringNumber: "124", number: 124},
		{stringNumber: "124211", number: 124211},
	}
	for _, item := range cases {
		num, err := parser.ParseNumber(item.stringNumber)
		assert.Equal(t, item.number, num)
		assert.Empty(t, err)
	}
}

func TestParseNumberWithInvalidNumber(t *testing.T) {
	t.Parallel()
	cases := []string{
		"124invalid",
		"_13",
	}
	for _, item := range cases {
		num, err := parser.ParseNumber(item)
		assert.Equal(t, uint(0), num)
		assert.NotEmpty(t, err)
	}
}
