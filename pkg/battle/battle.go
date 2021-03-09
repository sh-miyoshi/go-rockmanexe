package battle

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
)

// Init ...
func Init(plyr *player.Player) error {
	if err := fieldInit(); err != nil {
		return fmt.Errorf("Battle field init failed: %w", err)
	}

	return nil
}

// End ...
func End() {
	fieldEnd()
}

// Process ...
func Process() {}

// Draw ...
func Draw() {
	fieldDraw()
}
