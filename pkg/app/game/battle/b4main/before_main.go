package b4main

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	battleplayer "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	initShowDelay = 18
	showChipDelay = 10
)

type BeforeMain struct {
	msgInst      *titlemsg.TitleMsg
	count        int
	chipList     []chip.SelectParam
	paStartIndex int
	paEndIndex   int
	paID         int
	paShowTime   int
	waitTime     int
}

func New(selectChips []battleplayer.SelectChip) (*BeforeMain, error) {
	fname := config.ImagePath + "battle/msg_start.png"
	inst, err := titlemsg.New(fname, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize title message instance")
	}

	list := []chip.SelectParam{}
	for _, sc := range selectChips {
		c := chip.Get(sc.ID)
		list = append(list, chip.SelectParam{ID: c.ID, Code: sc.Code, Name: c.Name})
	}
	start, end, paID := chip.GetPAinList(list)

	return &BeforeMain{
		msgInst:      inst,
		chipList:     list,
		paStartIndex: start,
		paEndIndex:   end,
		paID:         paID,
		paShowTime:   (end-start)*showChipDelay + initShowDelay + 20,
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
		dxlib.DrawBox(0, 0, config.ScreenSize.X, config.ScreenSize.Y, 0x000000, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)

		draw.PAText(30, 50)

		for i, c := range b.chipList {
			if i < b.paStartIndex || i >= b.paEndIndex {
				draw.String(20, 100+i*40, 0xffffff, "%s", c.Name)
				draw.String(160, 105+i*40, 0xffffff, "%s", c.Code)
			} else {
				// Show PA data
				colors := []uint{
					dxlib.GetColor(248, 248, 176),
					dxlib.GetColor(248, 152, 0),
					dxlib.GetColor(56, 248, 144),
				}
				col := colors[(b.count/15)%len(colors)]
				if b.count > b.paShowTime {
					pa := chip.Get(b.paID)
					draw.String(20, 100+b.paStartIndex*40, col, "%s", pa.Name)
				} else if b.count > (i-b.paStartIndex)*showChipDelay+initShowDelay {
					draw.String(20, 100+i*40, col, "%s", c.Name)
					draw.String(160, 105+i*40, col, "%s", c.Code)
				}
			}
		}

		return
	}

	if b.msgInst != nil {
		b.msgInst.Draw()
	}
}

func (b *BeforeMain) Process() bool {
	b.count++

	if b.paID != -1 {
		if b.waitTime > 0 {
			b.waitTime--
			if b.waitTime == 0 {
				b.paID = -1
			}

			return false
		}

		if b.count == b.paShowTime {
			sound.On(resources.SEPACreated)
			b.waitTime = 60
		}
		if b.count < b.paShowTime-20 && b.count > initShowDelay && b.count%showChipDelay == 0 {
			sound.On(resources.SEPAPrepare)
		}

		return false
	}

	if b.msgInst != nil {
		return b.msgInst.Process()
	}
	return false
}
