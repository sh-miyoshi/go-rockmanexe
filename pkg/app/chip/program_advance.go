package chip

type paInfo struct {
	inputs []SelectParam
	paID   int
}

var (
	paList = []paInfo{
		{
			inputs: []SelectParam{
				{ID: IDSword, Code: ""},
				{ID: IDWideSword, Code: ""},
				{ID: IDLongSword, Code: ""},
			},
			paID: IDDreamSword,
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

func GetPAinList(chipList []SelectParam) (start, end int, paID int) {
	for i, c := range chipList {
		for _, pa := range paList {
			if expectChip(c, pa.inputs[0]) {
				ok := true
				for j := 1; j < len(pa.inputs); j++ {
					if i+j >= len(chipList) || !expectChip(chipList[i+j], pa.inputs[j]) {
						ok = false
						break
					}
				}
				if ok {
					start = i
					end = i + len(pa.inputs)
					paID = pa.paID
					return
				}
			}
		}
	}
	return -1, -1, -1
}

func expectChip(target, data SelectParam) bool {
	return target.ID == data.ID && (data.Code == "" || data.Code == target.Code)
}
