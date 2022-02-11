package chip

import "testing"

func TestGetPAinList(t *testing.T) {
	tt := []struct {
		ChipIDs    []int
		ExpectPAID int
	}{
		{
			ChipIDs:    []int{IDSword, IDWideSword, IDLongSword},
			ExpectPAID: IDDreamSword,
		},
		{
			ChipIDs:    []int{IDCannon, IDSword, IDWideSword, IDLongSword},
			ExpectPAID: IDDreamSword,
		},
		{
			ChipIDs:    []int{IDSword, IDWideSword, IDLongSword, IDCannon},
			ExpectPAID: IDDreamSword,
		},
		{
			ChipIDs:    []int{IDSword, IDWideSword, IDCannon},
			ExpectPAID: -1,
		},
		{
			ChipIDs:    []int{IDCannon, IDSword, IDWideSword},
			ExpectPAID: -1,
		},
		{
			ChipIDs:    []int{IDLongSword, IDSword, IDWideSword},
			ExpectPAID: -1,
		},
	}

	for i, tc := range tt {
		_, _, resID := GetPAinList(tc.ChipIDs)
		if resID != tc.ExpectPAID {
			t.Errorf("GetPAinList %d test failed. expect %v, but got %v", i, tc.ExpectPAID, resID)
		}
	}
}
