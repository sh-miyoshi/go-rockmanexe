package draw

import (
	"fmt"
)

var (
	objDrawInst  objectDraw
	animDrawInst animDraw
)

func Init() error {
	if err := objDrawInst.Init(); err != nil {
		return fmt.Errorf("object draw init failed: %w", err)
	}

	if err := animDrawInst.Init(); err != nil {
		return fmt.Errorf("anim draw init failed: %w", err)
	}

	return nil
}

func End() {
	objDrawInst.End()
	animDrawInst.End()
}

func Draw() {
	objDrawInst.Draw()
	animDrawInst.Draw()
}
