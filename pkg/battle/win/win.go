package win

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/titlemsg"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
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
	imgFrame   int32
	imgWinIcon int32
	count      int
	state      int
	deleteTime int
	winMsgInst *titlemsg.TitleMsg
)

func Init(gameTime int) error {
	state = stateMsg
	deleteTime = gameTime
	count = 0

	fname := common.ImagePath + "battle/result_frame.png"
	imgFrame = dxlib.LoadGraph(fname)
	if imgFrame == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/win_icon.png"
	imgWinIcon = dxlib.LoadGraph(fname)
	if imgWinIcon == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/msg_win.png"
	var err error
	winMsgInst, err = titlemsg.New(fname)

	return err
}

func End() {
	dxlib.DeleteGraph(imgFrame)
	dxlib.DeleteGraph(imgWinIcon)
	winMsgInst.End()
	winMsgInst = nil
}

func Process() bool {
	count++

	switch state {
	case stateMsg:
		if winMsgInst != nil && winMsgInst.Process() {
			stateChange(stateFrameIn)
			return false
		}
	case stateFrameIn:
		if count > 60 {
			stateChange(stateResult)
			return false
		}
	case stateResult:
		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			return true
		}
	}

	return false
}

func Draw() {
	switch state {
	case stateMsg:
		if winMsgInst != nil {
			winMsgInst.Draw()
		}
	case stateFrameIn:
		x := int32(count * 2)
		if x > 45 {
			x = 45
		}
		dxlib.DrawGraph(x, 30, imgFrame, dxlib.TRUE)
	case stateResult:
		dxlib.DrawGraph(45, 30, imgFrame, dxlib.TRUE)
		dxlib.DrawGraph(285, 180, imgWinIcon, dxlib.TRUE)
		draw.String(105, 230, 0xffffff, "Winner バッチ")
		showDeleteTime()
	}
}

func stateChange(nextState int) {
	logger.Info("Change battle result win state from %d to %d", state, nextState)
	if nextState < 0 || nextState >= stateMax {
		panic(fmt.Sprintf("Invalid next battle result win state: %d", nextState))
	}
	state = nextState
	count = 0
}

func showDeleteTime() {
	tm := deleteTime / 60
	if tm == 0 {
		tm = 1
	}

	min := tm / 60
	sec := tm % 60
	if min > 99 {
		min = 99
	}
	zero := 0
	draw.Number(300, 77, int32(min), draw.NumberOption{Padding: &zero, Length: 2})
	draw.String(333, 80, 0xffffff, "：")
	draw.Number(350, 77, int32(sec), draw.NumberOption{Padding: &zero, Length: 2})
}
