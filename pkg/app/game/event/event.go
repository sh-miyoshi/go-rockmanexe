package event

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

const (
	TypeChangeMapArea int = iota
)

const (
	ResultContinue int = iota
	ResultMapChange
	ResultEnd
)

type Scenario struct {
	Type int
	// TODO values
}

type Handler interface {
	Init(args string) error
	Draw()
	Process() (int, error)
}

var (
	handler   Handler
	scenarios []Scenario
	current   = 0
)

func SetScenarios(s []Scenario) {
	if len(s) == 0 {
		// 何もしない
		return
	}

	scenarios = append([]Scenario{}, s...)
	current = 0
	setHandler(scenarios[current])
}

func Draw() {
	if handler != nil {
		handler.Draw()
	}
}

func Process() (int, error) {
	if handler != nil {
		res, err := handler.Process()
		if err != nil {
			return ResultContinue, err
		}
		if res != ResultContinue {
			current++
			if current >= len(scenarios) {
				return ResultEnd, nil
			}
			setHandler(scenarios[current])
			return res, nil
		}
	}
	return ResultContinue, nil
}

func setHandler(s Scenario) {
	switch s.Type {
	case TypeChangeMapArea:
		handler = &MapChangeHandler{}
	default:
		common.SetError(fmt.Sprintf("scenario type %d is not implemented", s.Type))
	}
}
