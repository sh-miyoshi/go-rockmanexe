package chip

import "testing"

func TestGetPAinList(t *testing.T) {
	tt := []struct {
		ChipList   []SelectParam
		ExpectPAID int
	}{
		{
			ChipList: []SelectParam{
				{ID: IDSword, Code: ""},
				{ID: IDWideSword, Code: ""},
				{ID: IDLongSword, Code: ""},
			},
			ExpectPAID: IDDreamSword,
		},
		{
			ChipList: []SelectParam{
				{ID: IDCannon, Code: ""},
				{ID: IDSword, Code: ""},
				{ID: IDWideSword, Code: ""},
				{ID: IDLongSword, Code: ""},
			},
			ExpectPAID: IDDreamSword,
		},
		{
			ChipList: []SelectParam{
				{ID: IDSword, Code: ""},
				{ID: IDWideSword, Code: ""},
				{ID: IDLongSword, Code: ""},
				{ID: IDCannon, Code: ""},
			},
			ExpectPAID: IDDreamSword,
		},
		{
			ChipList: []SelectParam{
				{ID: IDSword, Code: ""},
				{ID: IDWideSword, Code: ""},
				{ID: IDCannon, Code: ""},
			},
			ExpectPAID: -1,
		},
		{
			ChipList: []SelectParam{
				{ID: IDCannon, Code: ""},
				{ID: IDSword, Code: ""},
				{ID: IDWideSword, Code: ""},
			},
			ExpectPAID: -1,
		},
		{
			ChipList: []SelectParam{
				{ID: IDLongSword, Code: ""},
				{ID: IDSword, Code: ""},
				{ID: IDWideSword, Code: ""},
			},
			ExpectPAID: -1,
		},
		{
			ChipList: []SelectParam{
				{ID: IDSword, Code: "s"},
				{ID: IDWideSword, Code: "s"},
				{ID: IDLongSword, Code: "s"},
			},
			ExpectPAID: IDDreamSword,
		},
	}

	for i, tc := range tt {
		_, _, resID := GetPAinList(tc.ChipList)
		if resID != tc.ExpectPAID {
			t.Errorf("GetPAinList %d test failed. expect %v, but got %v", i, tc.ExpectPAID, resID)
		}
	}
}
