package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/fps"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type menuPlayerStatus struct {
	playerInfo *player.Player
}

func playerStatusNew(plyr *player.Player) (*menuPlayerStatus, error) {
	res := menuPlayerStatus{
		playerInfo: plyr,
	}
	return &res, nil
}

func (r *menuPlayerStatus) End() {
}

func (r *menuPlayerStatus) Process() {
	if inputs.CheckKey(inputs.KeyCancel) == 1 {
		sound.On(resources.SECancel)
		stateChange(stateTop)
	}
}

func (r *menuPlayerStatus) Draw() {
	dxlib.DrawBox(60, 35, 400, 270, dxlib.GetColor(168, 192, 216), true)

	// get game count as seconds
	tm := r.playerInfo.PlayCount / uint(fps.FPS)
	if tm > 999*12*60 {
		tm = 999 * 12 * 60
	}
	tm /= 60 // change to minutes

	chipNum := player.FolderSize
	chipNum += len(r.playerInfo.BackPack)

	info := []struct {
		key   string
		value string
	}{
		{
			key:   "プレイ時間",
			value: fmt.Sprintf("%03d：%02d", tm/12, tm%12),
		},
		{
			key:   "バトルチップ",
			value: fmt.Sprintf("%d枚", chipNum),
		},
		{
			key:   "お金",
			value: fmt.Sprintf("%d ゼニー", r.playerInfo.Zenny),
		},
		{
			key:   "",
			value: fmt.Sprintf(""),
		},
		{
			key:   "ＨＰ",
			value: fmt.Sprintf("%d", r.playerInfo.HP),
		},
		{
			key:   "アタックレベル",
			value: fmt.Sprintf("%d", r.playerInfo.ShotPower),
		},
		{
			key:   "チャージ時間",
			value: fmt.Sprintf("%.01f秒", float64(r.playerInfo.ChargeTime)/float64(fps.FPS)),
		},
	}

	for i, row := range info {
		draw.String(80, 50+i*30, 0, row.key)
		draw.String(280, 50+i*30, 0, row.value)
	}
}
