package netbattle

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	netconn "github.com/sh-miyoshi/go-rockmanexe/pkg/app/newnetconn"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

var (
	InvalidChips = []int{
		chip.IDBoomerang1,
	}
)

func Init(plyr *player.Player) error {
	logger.Info("Init net battle data ...")
	return nil
}

func End() {
	netconn.GetInst().Disconnect()
}

func Process() error {
	return nil
}

func Draw() {
}
