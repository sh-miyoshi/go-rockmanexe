package navicustom

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateOpening int = iota
	stateMain
	stateRun
)

var (
	state int
	count int
)

func Init() error {
	state = stateOpening
	count = 0
	return nil
}

func Draw() {
	switch state {
	case stateOpening:
		// 起動時アニメーション
	case stateMain:
		// 実際にパーツを置いたりする
	case stateRun:
		// RUN
	}
}

func Process() {
	switch state {
	case stateOpening:
		if count > 30 || (count > 3 && inputs.CheckKey(inputs.KeyEnter) > 0) {
			stateChange(stateMain)
		}
	case stateMain:
		// TODO
	case stateRun:
		// TODO
	}

	count++
}

func stateChange(next int) {
	logger.Info("Change navu cutom state from %d to %d", state, next)
	state = next
	count = 0
}
