package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type menuInvalidChip struct {
	imgMsgFrame int
}

func invalidChipNew() (*menuInvalidChip, error) {
	res := menuInvalidChip{}

	fname := common.ImagePath + "menu/msg_frame.png"
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
		sound.On(sound.SECancel)
		stateChange(stateTop)
	}
}

func (i *menuInvalidChip) Draw() {
	dxlib.DrawBox(25, 45, 460, 200, dxlib.GetColor(168, 192, 216), true)
	draw.MessageText(35, 55, 0x000000, "使用できないチップ一覧")
	for i, cid := range netbattle.InvalidChips {
		draw.MessageText(35, 80+(i*30), 0x000000, fmt.Sprintf("・%s", chip.Get(cid).Name))
	}

	dxlib.DrawGraph(40, 205, i.imgMsgFrame, true)
	draw.MessageText(120, 220, 0x000000, "これらのチップはまだ通信対戦では使えな")
	draw.MessageText(120, 220+28, 0x000000, "いんだ。")
	draw.MessageText(120, 220+56, 0x000000, "チップフォルダを編集しよう")
}
