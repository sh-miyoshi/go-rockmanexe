package mapmove

import (
	"errors"
)

var (
	ErrGoBattle = errors.New("go to battle")
	ErrGoMenu   = errors.New("go to menu")
)

func Init() error {
	return nil
}

func End() {
}

func Draw() {
}

func Process() error {
	return nil
}
