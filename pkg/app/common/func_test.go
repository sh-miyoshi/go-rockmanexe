package common

import "testing"

func TestAbs(t *testing.T) {
	tt := []struct {
		input  int
		expect int
	}{
		{
			input:  5,
			expect: 5,
		},
		{
			input:  -5,
			expect: 5,
		},
		{
			input:  0,
			expect: 0,
		},
	}

	for _, tc := range tt {
		res := Abs(tc.input)
		if res != tc.expect {
			t.Errorf("Abs method expects %d, but got %d", tc.expect, res)
		}
	}
}
