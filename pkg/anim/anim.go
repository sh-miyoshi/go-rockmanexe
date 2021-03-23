package anim

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
)

const (
	TypeObject int = iota
	TypeEffect
)

type Param struct {
	PosX     int
	PosY     int
	AnimType int
}

// Anim ...
type Anim interface {
	Process() (bool, error)
	Draw()
	DamageProc(dm *damage.Damage)
	GetParam() Param
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

	for _, anim := range anims {
		pm := anim.GetParam()
		if dm := damage.Get(pm.PosX, pm.PosY); dm != nil {
			anim.DamageProc(dm)
			damage.Remove(dm.ID)
			// TODO if !pm.Penetrate delete(anims, id)
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
