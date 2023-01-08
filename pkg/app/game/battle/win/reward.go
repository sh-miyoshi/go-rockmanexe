package win

import (
	"math/rand"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
)

const (
	rewardTypeMoney int = iota
	rewardTypeChip
)

type rewardInfo struct {
	Type  int
	Name  string
	Image int
	Value interface{}
}

func getReward(all []rewardInfo) rewardInfo {
	if len(all) == 0 {
		common.SetError("no reward data")
		return rewardInfo{}
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
