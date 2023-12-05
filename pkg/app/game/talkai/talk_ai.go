package talkai

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const (
	stateInput = iota
	stateOutput
)

const (
	inputMaxLen = 80
)

var (
	state       int
	inputHandle int
	win         window.MessageWindow
)

func Init() {
	state = stateInput
	inputHandle = dxlib.MakeKeyInput(inputMaxLen, false, false, false, false, false)
	win.Init()
	background.Set(background.Type秋原町)
	win.SetMessage("", window.FaceTypeRockman)
	dxlib.SetActiveKeyInput(inputHandle)
	b := dxlib.GetColor(0, 0, 0)
	w := dxlib.GetColor(255, 255, 255)
	dxlib.SetKeyInputStringColor(b, b, w, b, b, w, b, b, b, b, b, w, w, b, b, b, b)
}

func End() {
	win.End()
	background.Unset()
	dxlib.DeleteKeyInput(inputHandle)
}

func Draw() {
	background.Draw()
	win.Draw()
	dxlib.DrawBox(45, 75, 430, 140, dxlib.GetColor(232, 184, 56), true)
	dxlib.DrawFormatString(50, 80, 0x000000, "質問を入力してね")
	dxlib.DrawBox(55, 100, 420, 130, 0xffffff, true)

	switch state {
	case stateInput:
		dxlib.DrawKeyInputString(65, 110, inputHandle, true)
	case stateOutput:
	}
}

func Process() bool {
	background.Process()

	switch state {
	case stateInput:
		if dxlib.CheckKeyInput(inputHandle) {
			str := fmt.Sprintf("%sって入力されたよ", inputString(inputHandle))
			win.SetMessage(str, window.FaceTypeRockman)
			state = stateOutput
		}
	case stateOutput:
		return win.Process()
	}
	return false
}

func inputString(handle int) string {
	buf := make([]byte, inputMaxLen)
	dxlib.GetKeyInputString(buf, inputHandle)
	slicedBuf := []byte{}
	for i := 0; i < inputMaxLen; i++ {
		if buf[i] == 0 {
			break
		}
		slicedBuf = append(slicedBuf, buf[i])
	}

	t := japanese.ShiftJIS.NewDecoder()
	str, _, _ := transform.Bytes(t, slicedBuf)
	return string(str)
}
