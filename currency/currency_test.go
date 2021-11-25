package currency_test

import (
	"fmt"
	"testing"

	"github.com/justteddy/wallet/currency"
	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	cases := []struct {
		cents    int
		expected string
	}{
		{
			cents:    0,
			expected: "0.00$",
		},
		{
			cents:    1,
			expected: "0.01$",
		},
		{
			cents:    10,
			expected: "0.10$",
		},
		{
			cents:    43,
			expected: "0.43$",
		},
		{
			cents:    90,
			expected: "0.90$",
		},
		{
			cents:    99,
			expected: "0.99$",
		},
		{
			cents:    100,
			expected: "1.00$",
		},
		{
			cents:    101,
			expected: "1.01$",
		},
		{
			cents:    110,
			expected: "1.10$",
		},
		{
			cents:    1199,
			expected: "11.99$",
		},
		{
			cents:    100_000,
			expected: "1000.00$",
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			assert.Equal(t, c.expected, currency.Format(c.cents))
		})
	}
}
