package opening

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	endCount = 90
)

var (
	count int32
)

func Init() {
	count = 0
}

func Process() bool {
	if config.Get().Debug.SkipBattleOpening {
		return true
	}

	if count == 0 {
		sound.On(sound.SEGoBattle)
	}

	count++
	return count > endCount
}

func Draw() {
	val := 255 * (endCount - count) / endCount
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, val)
	dxlib.DrawBox(0, 0, common.ScreenX, common.ScreenY, 0, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
}
