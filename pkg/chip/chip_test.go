package chip

import "testing"

func TestSelectable(t *testing.T) {
	tt := []struct {
		CurrentList []SelectParam
		Target      SelectParam
		Expect      bool
	}{
		{ // Case 0(Same name)
			CurrentList: []SelectParam{
				{Name: "name", Code: "a"},
				{Name: "name", Code: "b"},
			},
			Target: SelectParam{Name: "name", Code: "c"},
			Expect: true,
		},
		{ // Case 1(Same code)
			CurrentList: []SelectParam{
				{Name: "name1", Code: "a"},
				{Name: "name2", Code: "a"},
			},
			Target: SelectParam{Name: "name3", Code: "a"},
			Expect: true,
		},
		{ // Case 2(Same code with * in list)
			CurrentList: []SelectParam{
				{Name: "name1", Code: "*"},
				{Name: "name2", Code: "*"},
			},
			Target: SelectParam{Name: "name3", Code: "a"},
			Expect: true,
		},
		{ // Case 3(Same code with * in target)
			CurrentList: []SelectParam{
				{Name: "name1", Code: "a"},
				{Name: "name2", Code: "a"},
			},
			Target: SelectParam{Name: "name3", Code: "*"},
			Expect: true,
		},
		{ // Case 4(Same code with * in both)
			CurrentList: []SelectParam{
				{Name: "name1", Code: "a"},
				{Name: "name2", Code: "*"},
			},
			Target: SelectParam{Name: "name3", Code: "*"},
			Expect: true,
		},
		{ // Case 5(No selectable)
			CurrentList: []SelectParam{
				{Name: "name1", Code: "a"},
				{Name: "name2", Code: "a"},
			},
			Target: SelectParam{Name: "name1", Code: "b"},
			Expect: false,
		},
		{ // Case 6(Same name in list, and target is * code)
			CurrentList: []SelectParam{
				{Name: "name1", Code: "a"},
				{Name: "name1", Code: "b"},
			},
			Target: SelectParam{Name: "name2", Code: "*"},
			Expect: false,
		},
	}

	for i, tc := range tt {
		res := Selectable(tc.Target, tc.CurrentList)
		if res != tc.Expect {
			t.Errorf("Selectable %d test failed. expect %v, but got %v", i, tc.Expect, res)
		}
	}
}
