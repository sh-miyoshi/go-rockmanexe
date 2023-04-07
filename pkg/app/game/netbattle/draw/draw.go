package draw

import (
	"fmt"
)

var (
	objDrawInst objectDraw
)

func Init() error {
	if err := objDrawInst.Init(); err != nil {
		return fmt.Errorf("object draw init failed: %w", err)
	}

	return nil
}

func End() {
	objDrawInst.End()
}

func Draw() {
	objDrawInst.Draw()
}
