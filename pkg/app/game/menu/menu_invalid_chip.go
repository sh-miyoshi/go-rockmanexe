package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/stretchr/stew/slice"
)

type menuInvalidChip struct {
	playerInfo  *player.Player
	imgMsgFrame int
}

func invalidChipNew(plyr *player.Player) (*menuInvalidChip, error) {
	res := menuInvalidChip{
		playerInfo: plyr,
	}

	fname := common.ImagePath + "msg_frame.png"
	res.imgMsgFrame = dxlib.LoadGraph(fname)
	if res.imgMsgFrame == -1 {
		return nil, fmt.Errorf("failed to load menu message frame image %s", fname)
	}

	return &res, nil
}

func (i *menuInvalidChip) End() {
	dxlib.DeleteGraph(i.imgMsgFrame)
}

func (i *menuInvalidChip) Process() {
	if inputs.CheckKey(inputs.KeyCancel) == 1 || inputs.CheckKey(inputs.KeyEnter) == 1 {
		sound.On(resources.SECancel)
		stateChange(stateTop)
	}
}

func (i *menuInvalidChip) Draw() {
	dxlib.DrawBox(25, 45, 460, 200, dxlib.GetColor(168, 192, 216), true)
	draw.MessageText(35, 55, 0x000000, "使用できないチップ一覧")
	for i, name := range i.invalidChips() {
		draw.MessageText(35, 80+(i*30), 0x000000, fmt.Sprintf("・%s", name))
	}

	dxlib.DrawGraph(40, 205, i.imgMsgFrame, true)
	draw.MessageText(120, 220, 0x000000, "これらのチップはまだ通信対戦では使えな")
	draw.MessageText(120, 220+28, 0x000000, "いんだ。")
	draw.MessageText(120, 220+56, 0x000000, "チップフォルダを編集しよう")
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
