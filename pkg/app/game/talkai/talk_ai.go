package talkai

import "github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"

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
}

func End() {
	win.End()
}

func Draw() {
	win.Draw()
}

func Process() bool {
	switch state {
	case stateInput:
		// TODO
	case stateOutput:
		return win.Process()
	}
	return false
}
