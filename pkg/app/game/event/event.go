package event

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
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
	Init(values []byte) error
	End()
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
	if err := setHandler(scenarios[current]); err != nil {
		common.SetError(fmt.Sprintf("failed to set handler: %v", err))
	}
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
			handler.End()
			return ResultContinue, err
		}
		if end {
			handler.End()

			current++
			if current >= len(scenarios) {
				return ResultEnd, nil
			}
			if err := setHandler(scenarios[current]); err != nil {
				return ResultContinue, err
			}
		}
	}
	return resultCode, nil
}

func setHandler(s Scenario) error {
	switch s.Type {
	case TypeChangeMapArea:
		handler = &MapChangeHandler{}
	case TypeEnd:
		handler = &EndHandler{}
	case TypeMessage:
		handler = &MessageHandler{}
	default:
		return fmt.Errorf("scenario type %d is not implemented", s.Type)
	}
	return handler.Init(s.Values)
}
