package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/stretchr/stew/slice"
)

type menuInvalidChip struct {
	playerInfo *player.Player
	win        window.MessageWindow
}

func invalidChipNew(plyr *player.Player) (*menuInvalidChip, error) {
	res := menuInvalidChip{
		playerInfo: plyr,
	}

	if err := res.win.Init(); err != nil {
		return nil, err
	}
	msg := "これらのチップはまだ通信対戦では使えないんだ。チップフォルダを編集しよう"
	res.win.SetMessage(msg, window.FaceTypeRockman)

	return &res, nil
}

func (i *menuInvalidChip) End() {
	i.win.End()
}

func (i *menuInvalidChip) Process() bool {
	if i.win.Process() {
		sound.On(resources.SECancel)
		return true
	}
	return false
}

func (i *menuInvalidChip) Draw() {
	dxlib.DrawBox(25, 45, 460, 200, dxlib.GetColor(168, 192, 216), true)
	draw.MessageText(35, 55, 0x000000, "使用できないチップ一覧")
	for i, name := range i.invalidChips() {
		draw.MessageText(35, 80+(i*30), 0x000000, fmt.Sprintf("・%s", name))
	}

	i.win.Draw()
}

func (i *menuInvalidChip) GetResult() Result {
	return ResultContinue
}

func (i *menuInvalidChip) invalidChips() []string {
	res := []string{}
	for _, c := range i.playerInfo.ChipFolder {
		if !slice.Contains(netbattle.ValidChips, c.ID) {
			res = append(res, chip.Get(c.ID).Name)
			if len(res) >= 3 {
				res = append(res, "など・・・")
				return res
			}
		}
	}
	return res
}
