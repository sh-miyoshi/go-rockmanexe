package chip

/*
やること
chipListを入手
n~m番目がどのチップになるか返す

プログラムアドバンス一覧を管理
*/

type paInfo struct {
	inputIDs []int
	paID     int
}

var (
	paList = []paInfo{
		{
			inputIDs: []int{IDSword, IDWideSword, IDLongSword},
			paID:     IDDreamSword,
		},
	}
)

func setPAData() {
	// Set program advance to chip list
	chipData = append(chipData, Chip{
		ID:               IDDreamSword,
		Name:             "ドリームソード",
		Power:            400,
		Type:             TypeSword,
		PlayerAct:        4,
		IsProgramAdvance: true,
	})
}

func GetPAinList(chipIDs []int) (start, end int, paID int) {
	for i, cid := range chipIDs {
		for _, pa := range paList {
			if pa.inputIDs[0] == cid {
				ok := true
				for j := 1; j < len(pa.inputIDs); j++ {
					if i+j >= len(chipIDs) || pa.inputIDs[j] != chipIDs[i+j] {
						ok = false
						break
					}
				}
				if ok {
					start = i
					end = i + len(pa.inputIDs)
					paID = pa.paID
					return
				}
			}
		}
	}
	return -1, -1, -1
}
