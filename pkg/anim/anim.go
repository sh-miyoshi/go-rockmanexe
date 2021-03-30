package anim

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	TypeEffect int = iota
	TypeObject
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
	anims         = map[string]Anim{}
	sortedAnimIDs = []string{}
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
			for i, sid := range sortedAnimIDs {
				if sid == id {
					sortedAnimIDs = append(sortedAnimIDs[:i], sortedAnimIDs[i+1:]...)
				}
			}
		}
	}

	// Damage Process
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
	for _, id := range sortedAnimIDs {
		anims[id].Draw()
	}
}

// New ...
func New(anim Anim) string {
	id := uuid.New().String()
	pm := anim.GetParam()
	index := pm.AnimType*100 + pm.PosY*6 + pm.PosX

	if len(sortedAnimIDs) == 0 {
		sortedAnimIDs = append(sortedAnimIDs, id)
	} else {
		set := false
		for i, sid := range sortedAnimIDs {
			spm := anims[sid].GetParam()
			sindex := spm.AnimType*100 + spm.PosY*6 + spm.PosX
			if index > sindex {
				tmp := append([]string{}, sortedAnimIDs[i:]...)
				sortedAnimIDs = append(sortedAnimIDs[:i], id)
				sortedAnimIDs = append(sortedAnimIDs, tmp...)
				set = true
				break
			}
		}
		if !set {
			sortedAnimIDs = append(sortedAnimIDs, id)
		}
	}

	anims[id] = anim

	logger.Debug("added anim %s with %+v", id, pm)
	return id
}

// IsProcessing ...
func IsProcessing(id string) bool {
	_, exists := anims[id]
	return exists
}
