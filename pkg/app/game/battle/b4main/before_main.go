package b4main

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
)

type BeforeMain struct {
	msgInst *titlemsg.TitleMsg
}

func New(selectChips []int) (*BeforeMain, error) {
	fname := common.ImagePath + "battle/msg_start.png"
	inst, err := titlemsg.New(fname, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize title message instance: %w", err)
	}

	return &BeforeMain{
		msgInst: inst,
	}, nil
}

func (b *BeforeMain) End() {
	if b.msgInst != nil {
		b.msgInst.End()
	}
}

func (b *BeforeMain) Draw() {
	if b.msgInst != nil {
		b.msgInst.Draw()
	}
}

func (b *BeforeMain) Process() bool {
	if b.msgInst != nil {
		return b.msgInst.Process()
	}
	return false
}
