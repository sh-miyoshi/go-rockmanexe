package event

import "fmt"

const (
	TypeChangeMapArea int = iota
)

type Handler interface {
	Init(args string) error
	Draw()
	Process() error
}

var (
	handler Handler
)

func Set(eventType int, args string) {
	switch eventType {
	case TypeChangeMapArea:
		// TODO handler = &mapchange.Handler{}
	default:
		panic(fmt.Sprintf("invalid event type %d was specified", eventType))
	}
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
