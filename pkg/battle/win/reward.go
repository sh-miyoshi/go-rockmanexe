package win

import (
	"math/rand"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

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

func rewardProc(data rewardInfo, plyr *player.Player) {
	switch data.Type {
	case rewardTypeMoney:
		plyr.UpdateMoney(data.Value.(int))
	case rewardTypeChip:
		c := chip.GetByName(data.Name)
		plyr.AddChip(c.ID, data.Value.(string))
	}
}
