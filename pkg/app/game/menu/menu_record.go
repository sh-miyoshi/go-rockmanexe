package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type menuRecord struct {
	playerInfo *player.Player
}

func recordNew(plyr *player.Player) (*menuRecord, error) {
	res := menuRecord{
		playerInfo: plyr,
	}
	return &res, nil
}

func (r *menuRecord) End() {
}

func (r *menuRecord) Process() {
	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		sound.On(resources.SECancel)
		stateChange(stateTop)
	}
}

func (r *menuRecord) Draw() {
	// get game count as seconds (FPS: 60)
	tm := r.playerInfo.PlayCount / 60
	if tm > 999*12*60 {
		tm = 999 * 12 * 60
	}
	tm /= 60 // change to minutes

	chipNum := player.FolderSize
	chipNum += len(r.playerInfo.BackPack)

	draw.String(80, 50, 0, "プレイ時間                 %03d：%02d", tm/12, tm%12)
	draw.String(80, 90, 0, "バトルチップ              %6d枚", chipNum)
	draw.String(80, 130, 0, "お金                    %7d ゼニー", r.playerInfo.Zenny)
}
