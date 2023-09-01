package menu

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/list"
	"github.com/stretchr/stew/slice"
)

const (
	topSelectChipFolder int = iota
	topSelectGoBattle
	topSelectRecord
	topSelectNetBattle
)

type menuTop struct {
	playerInfo *player.Player
	itemList   list.ItemList
}

func topNew(plyr *player.Player) (*menuTop, error) {
	res := &menuTop{
		playerInfo: plyr,
	}
	res.itemList.SetList([]string{
		"チップフォルダ",
		"バトル",
		"戦績",
		"ネット対戦",
	})

	return res, nil
}

func (t *menuTop) End() {
}

func (t *menuTop) Process() {
	if config.Get().Debug.EnableDevFeature {
		if inputs.CheckKey(inputs.KeyLButton) == 1 {
			sound.On(resources.SEMenuEnter)
			stateChange(stateDevFeature)
			return
		}
	}

	sel := t.itemList.Process()
	if sel != -1 {
		sound.On(resources.SEMenuEnter)
		switch sel {
		case topSelectChipFolder:
			stateChange(stateChipFolder)
		case topSelectGoBattle:
			stateChange(stateGoBattle)
		case topSelectRecord:
			stateChange(stateRecord)
		case topSelectNetBattle:
			if t.haveInvalidChip() {
				stateChange(stateInvalidChip)
			} else {
				stateChange(stateNetBattle)
			}
		}
	}
}

func (t *menuTop) Draw() {
	dxlib.DrawBox(20, 30, 230, 300, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(30, 40, 210, len(t.itemList.GetList())*35+50, dxlib.GetColor(16, 80, 104), true)

	for i, msg := range t.itemList.GetList() {
		draw.String(65, 50+i*35, 0xffffff, msg)
	}

	const s = 2
	y := 50 + t.itemList.GetPointer()*35
	dxlib.DrawTriangle(40, y+s, 40+18-s*2, y+10, 40, y+20-s, 0xffffff, true)

	// Show description
	dxlib.DrawBox(255, 55, 445, 285, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(275, 38, 425, 55, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawTriangle(255, 55, 275, 38, 275, 55, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawTriangle(425, 55, 425, 38, 445, 55, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(260, 60, 440, 280, dxlib.GetColor(16, 80, 104), true)
	draw.String(280, 40, 0xffffff, "Description")

	switch t.itemList.GetPointer() {
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

	if config.Get().Debug.EnableDevFeature {
		draw.String(50, 220, 0x000000, "L-btn: Debug機能")
	}
}

func (t *menuTop) haveInvalidChip() bool {
	for _, c := range t.playerInfo.ChipFolder {
		if slice.Contains(netbattle.InvalidChips, c.ID) {
			return true
		}
	}
	return false
}
