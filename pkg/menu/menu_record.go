package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
)

func recordInit() error {
	return nil
}

func recordEnd() {
}

func recordProcess() {
	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		sound.On(sound.SECancel)
		stateChange(stateTop)
	}
}

func recordDraw() {
	draw.String(common.ScreenX/2-20, common.ScreenY/2-20, 0, "未実装")
}
