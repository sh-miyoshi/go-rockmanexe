package talkai

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/background"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
)

const (
	stateInput = iota
	stateOutput
)

var (
	state int
	win   window.MessageWindow
)

func Init() {
	state = stateInput
	win.Init()
	background.Set(background.Type秋原町)
	win.SetMessage("", window.FaceTypeRockman)
}

func End() {
	win.End()
	background.Unset()
}

func Draw() {
	background.Draw()
	win.Draw()
}

func Process() bool {
	background.Process()

	switch state {
	case stateInput:
		// TODO
	case stateOutput:
		return win.Process()
	}
	return false
}
