package event

import (
	"errors"
	"fmt"
)

const (
	TypeChangeMapArea int = iota
)

type Handler interface {
	Init(args string) error
	Draw()
	Process() error
}

var (
	ErrMapChange = errors.New("change map")
)

var (
	handler Handler
)

func Set(eventType int, args string) error {
	switch eventType {
	case TypeChangeMapArea:
		handler = &MapChangeHandler{}
	default:
		return fmt.Errorf("invalid event type %d was specified", eventType)
	}
	return handler.Init(args)
}

func Draw() {
	if handler != nil {
		handler.Draw()
	}
}

func Process() error {
	if handler != nil {
		return handler.Process()
	}
	return nil
}
