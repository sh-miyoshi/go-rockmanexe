package menu

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	topSelectChipFolder int = iota
	topSelectGoBattle
	topSelectRecord
	topSelectNetBattle

	topSelectMax
)

var (
	topPointer = 0
)

func topInit() error {
	return nil
}

func topEnd() {
}

func topProcess() {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		sound.On(sound.SEMenuEnter)
		switch topPointer {
		case topSelectChipFolder:
			stateChange(stateChipFolder)
		case topSelectGoBattle:
			stateChange(stateGoBattle)
		case topSelectRecord:
			stateChange(stateRecord)
		case topSelectNetBattle:
			stateChange(stateNetBattle)
		}
	} else {
		if inputs.CheckKey(inputs.KeyUp) == 1 {
			if topPointer > 0 {
				sound.On(sound.SECursorMove)
				topPointer--
			}
		} else if inputs.CheckKey(inputs.KeyDown) == 1 {
			if topPointer < topSelectMax-1 {
				sound.On(sound.SECursorMove)
				topPointer++
			}
		}
	}
}

func topDraw() {
	msgs := [topSelectMax]string{
		"チップフォルダ",
		"バトル",
		"戦績",
		"ネット対戦",
	}

	dxlib.DrawBox(20, 30, 230, 300, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawBox(30, 40, 210, int32(len(msgs)*35)+50, dxlib.GetColor(16, 80, 104), dxlib.TRUE)

	for i, msg := range msgs {
		draw.String(65, 50+int32(i)*35, 0xffffff, msg)
	}

	const s = 2
	y := int32(50 + topPointer*35)
	dxlib.DrawTriangle(40, y+s, 40+18-s*2, y+10, 40, y+20-s, 0xffffff, dxlib.TRUE)

	// Show description
	dxlib.DrawBox(255, 55, 445, 285, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawBox(275, 38, 425, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawTriangle(255, 55, 275, 38, 275, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawTriangle(425, 55, 425, 38, 445, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawBox(260, 60, 440, 280, dxlib.GetColor(16, 80, 104), dxlib.TRUE)
	draw.String(280, 40, 0xffffff, "Description")

	switch topPointer {
	case topSelectChipFolder:
		draw.String(270, 70, 0xffffff, "チップフォルダを閲覧し")
		draw.String(270, 100, 0xffffff, "ます")
	case topSelectGoBattle:
		draw.String(270, 70, 0xffffff, "ウィルスバスティングを")
		draw.String(270, 100, 0xffffff, "行います")
	case topSelectRecord:
		draw.String(270, 70, 0xffffff, "今までの戦績を確認しま")
		draw.String(270, 100, 0xffffff, "す")
	case topSelectNetBattle:
		draw.String(270, 70, 0xffffff, "インターネットを経由し")
		draw.String(270, 100, 0xffffff, "て対戦します")
	}
}
