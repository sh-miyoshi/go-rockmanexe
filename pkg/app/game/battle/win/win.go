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

var (
	imgFrame   int
	imgZenny   int
	count      int
	state      int
	winMsgInst *titlemsg.TitleMsg
)

func Init(args reward.WinArg, plyr *player.Player) error {
	state = stateMsg
	count = 0

	fname := config.ImagePath + "battle/result_frame.png"
	imgFrame = dxlib.LoadGraph(fname)
	if imgFrame == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/zenny.png"
	imgZenny = dxlib.LoadGraph(fname)
	if imgZenny == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	fname = config.ImagePath + "battle/msg_win.png"
	var err error
	winMsgInst, err = titlemsg.New(fname, 0)

	if err := sound.BGMPlay(sound.BGMWin); err != nil {
		return errors.Wrap(err, "failed to play bgm")
	}

	reward.SetToPlayer(args, plyr)
	return err
}

func End() {
	sound.SEClear()
	dxlib.DeleteGraph(imgFrame)
	dxlib.DeleteGraph(imgZenny)
	if winMsgInst != nil {
		winMsgInst.End()
		winMsgInst = nil
	}
	state = stateMsg
}

func Update() bool {
	count++

	switch state {
	case stateMsg:
		if winMsgInst != nil && winMsgInst.Update() {
			stateChange(stateFrameIn)
			return false
		}
	case stateFrameIn:
		if count > 60 {
			stateChange(stateResult)
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

func Draw() {
	baseX := config.ScreenSize.X/2 - 195
	baseY := config.ScreenSize.Y/2 - 130

	switch state {
	case stateMsg:
		if winMsgInst != nil {
			winMsgInst.Draw()
		}
	case stateFrameIn:
		x := count * baseX / 60
		if x > baseX {
			x = baseX
		}
		dxlib.DrawGraph(x, baseY, imgFrame, true)
	case stateResult:
		dxlib.DrawGraph(baseX, baseY, imgFrame, true)
		pm := reward.GetParam()
		switch pm.Type {
		case reward.TypeMoney:
			dxlib.DrawGraph(baseX+227, baseY+144, imgZenny, true)
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

func stateChange(nextState int) {
	logger.Info("Change battle result win state from %d to %d", state, nextState)
	if nextState < 0 || nextState >= stateMax {
		system.SetError(fmt.Sprintf("Invalid next battle result win state: %d", nextState))
	}
	state = nextState
	count = 0
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
