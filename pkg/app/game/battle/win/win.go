package win

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	chipimage "github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip/image"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/win/reward"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateMsg int = iota
	stateFrameIn
	stateResult

	stateMax
)

type WinManager struct {
	imgFrame   int
	imgZenny   int
	count      int
	state      int
	winMsgInst *titlemsg.TitleMsg
}

func New(args reward.WinArg, plyr *player.Player) (*WinManager, error) {
	res := &WinManager{
		state: stateMsg,
		count: 0,
	}

	fname := config.ImagePath + "battle/result_frame.png"
	res.imgFrame = dxlib.LoadGraph(fname)
	if res.imgFrame == -1 {
		return nil, errors.Newf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/zenny.png"
	res.imgZenny = dxlib.LoadGraph(fname)
	if res.imgZenny == -1 {
		return nil, errors.Newf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/msg_win.png"
	var err error
	res.winMsgInst, err = titlemsg.New(fname, 0)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create title msg")
	}

	if err := sound.BGMPlay(sound.BGMWin); err != nil {
		return nil, errors.Wrap(err, "failed to play bgm")
	}

	reward.SetToPlayer(args, plyr)
	return res, nil
}

func (m *WinManager) End() {
	sound.SEClear()
	dxlib.DeleteGraph(m.imgFrame)
	dxlib.DeleteGraph(m.imgZenny)
	if m.winMsgInst != nil {
		m.winMsgInst.End()
	}
}

func (m *WinManager) Update() bool {
	m.count++

	switch m.state {
	case stateMsg:
		if m.winMsgInst != nil && m.winMsgInst.Update() {
			m.stateChange(stateFrameIn)
			return false
		}
	case stateFrameIn:
		if m.count > 60 {
			m.stateChange(stateResult)
			sound.On(resources.SEGotItem)
			return false
		}
	case stateResult:
		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			sound.On(resources.SESelect)
			return true
		}
	}

	return false
}

func (m *WinManager) Draw() {
	baseX := config.ScreenSize.X/2 - 195
	baseY := config.ScreenSize.Y/2 - 130

	switch m.state {
	case stateMsg:
		if m.winMsgInst != nil {
			m.winMsgInst.Draw()
		}
	case stateFrameIn:
		x := m.count * baseX / 60
		if x > baseX {
			x = baseX
		}
		dxlib.DrawGraph(x, baseY, m.imgFrame, true)
	case stateResult:
		dxlib.DrawGraph(baseX, baseY, m.imgFrame, true)
		pm := reward.GetParam()
		switch pm.Type {
		case reward.TypeMoney:
			dxlib.DrawGraph(baseX+227, baseY+144, m.imgZenny, true)
		case reward.TypeChip:
			chipInfo := chip.GetByName(pm.Name)
			img := chipimage.GetDetail(chipInfo.ID)
			dxlib.DrawGraph(baseX+227, baseY+144, img, true)
			c := strings.ToUpper(pm.Code)
			draw.String(baseX+195, baseY+200, 0xffffff, c)
		}
		draw.String(baseX+60, baseY+200, 0xffffff, pm.Name)
		showDeleteTime(pm.DeleteTimeSec, baseX, baseY)
		draw.Number(baseX+315, baseY+95, pm.BustingLevel)
	}
}

func (m *WinManager) stateChange(nextState int) {
	logger.Info("Change battle result win state from %d to %d", m.state, nextState)
	if nextState < 0 || nextState >= stateMax {
		system.SetError(fmt.Sprintf("Invalid next battle result win state: %d", nextState))
	}
	m.state = nextState
	m.count = 0
}

func showDeleteTime(deleteTimeSec int, baseX, baseY int) {
	tm := deleteTimeSec

	min := tm / 60
	sec := tm % 60
	if min > 99 {
		min = 99
	}
	zero := 0
	draw.Number(baseX+255, baseY+47, min, draw.NumberOption{Padding: &zero, Length: 2})
	draw.String(baseX+288, baseY+50, 0xffffff, "：")
	draw.Number(baseX+305, baseY+47, sec, draw.NumberOption{Padding: &zero, Length: 2})
}
