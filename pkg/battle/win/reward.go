package win

import "math/rand"

const (
	rewardTypeMoney int = iota
	rewardTypeChip
)

type rewardInfo struct {
	Type  int
	Name  string
	Image int32
	Value interface{}
}

func getReward(all []rewardInfo) rewardInfo {
	if len(all) == 0 {
		panic("no reward data")
	}

	// TODO 重みづけ
	n := rand.Intn(len(all))
	return all[n]
}

func rewardProc(data rewardInfo) {
	// TODO implement this
}
