package event

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	TypeChangeMapArea int = iota
	TypeEnd
	TypeMessage
)

const (
	ResultContinue int = iota
	ResultMapChange
	ResultEnd
)

type Scenario struct {
	Type   int
	Values []byte
}

type Handler interface {
	Draw()
	Process() (bool, error)
}

var (
	handler      Handler
	scenarios    []Scenario
	current      = 0
	storedValues []byte
	resultCode   = ResultContinue
)

func SetScenarios(s []Scenario) {
	if len(s) == 0 {
		// 何もしない
		return
	}

	scenarios = append([]Scenario{}, s...)

	// 最後にEndHandlerを追加する
	scenarios = append(scenarios, Scenario{Type: TypeEnd})

	current = 0
	setHandler(scenarios[current])
	resultCode = ResultContinue
}

func GetStoredValues() []byte {
	return storedValues
}

func Draw() {
	if handler != nil {
		handler.Draw()
	}
}

func Process() (int, error) {
	if handler != nil {
		end, err := handler.Process()
		if err != nil {
			return ResultContinue, err
		}
		if end {
			current++
			if current >= len(scenarios) {
				return ResultEnd, nil
			}
			setHandler(scenarios[current])
		}
	}
	return resultCode, nil
}

func setHandler(s Scenario) {
	switch s.Type {
	case TypeChangeMapArea:
		var args MapChangeArgs
		args.Unmarshal(s.Values)
		handler = &MapChangeHandler{args: args}
		logger.Info("set map change scenario with %+v", args)
	case TypeEnd:
		handler = &EndHandler{}
		logger.Info("set end scenario")
	case TypeMessage:
		handler = &MessageHandler{Message: string(s.Values)}
		logger.Info("set message scenario with %s", string(s.Values))
	default:
		common.SetError(fmt.Sprintf("scenario type %d is not implemented", s.Type))
	}
}
