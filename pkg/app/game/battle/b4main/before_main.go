package b4main

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type BeforeMain struct {
	msgInst      *titlemsg.TitleMsg
	count        int
	chipList     []chip.SelectParam
	paStartIndex int
	paEndIndex   int
	paID         int
}

func New(selectChips []player.ChipInfo) (*BeforeMain, error) {
	fname := common.ImagePath + "battle/msg_start.png"
	inst, err := titlemsg.New(fname, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize title message instance: %w", err)
	}

	list := []chip.SelectParam{}
	for _, sc := range selectChips {
		c := chip.Get(sc.ID)
		list = append(list, chip.SelectParam{ID: c.ID, Code: sc.Code, Name: c.Name})
	}
	start, end, paID := chip.GetPAinList(list)

	/*
		ある時モード
			PA以外は描画しておく(白)
			PAは元チップをタイミングに合わせて描画(レインボー&SEあり)
			全部描画終わったらPAを描画(レインボー&SEあり)
		title msgに移る
	*/

	return &BeforeMain{
		msgInst:      inst,
		chipList:     list,
		paStartIndex: start,
		paEndIndex:   end,
		paID:         paID,
	}, nil
}

func (b *BeforeMain) End() {
	if b.msgInst != nil {
		b.msgInst.End()
	}
}

func (b *BeforeMain) Draw() {
	if b.paID != -1 {
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0x000000, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)

		for i, c := range b.chipList {
			draw.String(20, 100+i*40, 0xffffff, "%s", c.Name)
			draw.String(160, 105+i*40, 0xffffff, "%s", c.Code)
		}
		// TODO

		return
	}

	if b.msgInst != nil {
		b.msgInst.Draw()
	}
}

func (b *BeforeMain) Process() bool {
	b.count++

	if b.paID != -1 {
		// TODO
		return false
	}

	if b.msgInst != nil {
		return b.msgInst.Process()
	}
	return false
}
