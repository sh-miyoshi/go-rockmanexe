package ja

import "testing"

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
		{
			input:    "改行\nテスト",
			splitNum: 100,
			expect: []string{
				"改行",
				"テスト",
			},
		},
	}

	for _, tc := range tt {
		res := SplitMsg(tc.input, tc.splitNum)
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
