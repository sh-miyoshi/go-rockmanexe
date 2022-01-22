package mapmove

import (
	"errors"

	"github.com/sh-miyoshi/dxlib"
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
	dxlib.DrawString(0, 0, "Map Move", 0xffffff)
}

func Process() error {
	return nil
}
