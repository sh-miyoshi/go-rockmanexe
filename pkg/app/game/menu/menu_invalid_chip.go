package menu

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

type menuInvalidChip struct {
	win window.MessageWindow
}

func invalidChipNew() (*menuInvalidChip, error) {
	res := menuInvalidChip{}

	msg := "これらのチップはまだ通信対戦では使えないんだ。チップフォルダを編集しよう"
	var err error
	res.win, err = window.New(msg)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (i *menuInvalidChip) End() {
	i.win.End()
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
	for i, cid := range netbattle.InvalidChips {
		draw.MessageText(35, 80+(i*30), 0x000000, fmt.Sprintf("・%s", chip.Get(cid).Name))
	}

	i.win.Draw()
}
