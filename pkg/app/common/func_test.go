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

func TestSplitMsg(t *testing.T) {
	tt := []struct {
		input    string
		splitNum int
		expect   []string
	}{
		{
			input:    "hello,world",
			splitNum: 3,
			expect: []string{
				"hel",
				"lo,",
				"wor",
				"ld",
			},
		},
		{
			input:    "日本語テスト",
			splitNum: 3,
			expect: []string{
				"日本語",
				"テスト",
			},
		},
	}

	for _, tc := range tt {
		res := SplitJAMsg(tc.input, tc.splitNum)
		if len(res) != len(tc.expect) {
			t.Errorf("SplitMsg method expects %v, but got %v", tc.expect, res)
		} else {
			for i := 0; i < len(res); i++ {
				if res[i] != tc.expect[i] {
					t.Errorf("SplitMsg method expects %v, but got %v", tc.expect, res)
					break
				}
			}
		}
	}
}
