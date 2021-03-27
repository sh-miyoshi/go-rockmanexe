package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

func goBattleInit() error {
	return nil
}

func goBattleEnd() {
}

func goBattleProcess() {
	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		stateChange(stateTop)
	}
}

func goBattleDraw() {
	draw.String(common.ScreenX/2-20, common.ScreenY/2-20, 0, "未実装")
}
