package opening

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	endCount = 90
)

var (
	count int
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
	dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0, true)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
}
