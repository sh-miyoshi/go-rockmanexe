package win

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
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

const (
	msgDelay = 4
)

var (
	imgFrame   int32
	imgMsg     []int32
	imgWinIcon int32
	count      int
	state      int
	deleteTime int
)

func Init(gameTime int) error {
	state = stateMsg
	deleteTime = gameTime
	count = 0

	fname := common.ImagePath + "battle/result_frame.png"
	imgFrame = dxlib.LoadGraph(fname)
	if imgFrame == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
	}

	imgMsg = make([]int32, 3)
	fname = common.ImagePath + "battle/msg_win.png"
	if res := dxlib.LoadDivGraph(fname, 3, 1, 3, 272, 32, imgMsg); res == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
	}

	fname = common.ImagePath + "battle/win_icon.png"
	imgWinIcon = dxlib.LoadGraph(fname)
	if imgWinIcon == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgFrame)
	dxlib.DeleteGraph(imgWinIcon)
	for _, img := range imgMsg {
		dxlib.DeleteGraph(img)
	}
	imgMsg = []int32{}
}

func Process() bool {
	count++

	switch state {
	case stateMsg:
		if count >= len(imgMsg)*msgDelay+20 {
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
	if len(imgMsg) == 0 {
		// Waiting initialize
		return
	}

	switch state {
	case stateMsg:
		drawMsg()
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

func drawMsg() {
	imgNo := count / msgDelay
	if imgNo >= len(imgMsg) {
		imgNo = len(imgMsg) - 1
	}
	dxlib.DrawGraph(105, 125, imgMsg[imgNo], dxlib.TRUE)
}
