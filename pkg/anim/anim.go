package anim

import (
	"fmt"

	"github.com/google/uuid"
)

// Anim ...
type Anim interface {
	Process() (bool, error)
	Draw()
}

var (
	anims = map[string]Anim{}
)

// MgrProcess ...
func MgrProcess() error {
	for id, anim := range anims {
		end, err := anim.Process()
		if err != nil {
			return fmt.Errorf("Anim process failed: %w", err)
		}

		if end {
			delete(anims, id)
		}
	}
	return nil
}

// MgrDraw ...
func MgrDraw() {
	for _, anim := range anims {
		anim.Draw()
	}
}

// New ...
func New(anim Anim) string {
	id := uuid.New().String()
	anims[id] = anim
	return id
}

// IsProcessing ...
func IsProcessing(id string) bool {
	_, exists := anims[id]
	return exists
}
