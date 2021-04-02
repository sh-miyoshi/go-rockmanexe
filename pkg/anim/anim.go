package anim

import (
	"fmt"
	"sort"

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
	anims         = map[string]Anim{}
	sortedAnimIDs = []string{}
)

// MgrProcess ...
func MgrProcess(enableDamage bool) error {
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
	if enableDamage {
		for _, anim := range anims {
			pm := anim.GetParam()
			if dm := damage.Get(pm.PosX, pm.PosY); dm != nil {
				anim.DamageProc(dm)
				damage.Remove(dm.ID)
				// TODO if !pm.Penetrate delete(anims, id)
			}
		}
	}

	sortAnim()

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
	anims[id] = anim
	sortAnim()
	return id
}

// IsProcessing ...
func IsProcessing(id string) bool {
	_, exists := anims[id]
	return exists
}

func Cleanup() {
	anims = map[string]Anim{}
	sortedAnimIDs = []string{}
}

func sortAnim() {
	type sortParam struct {
		Index int
		ID    string
	}
	sortAnims := []sortParam{}
	for id, anim := range anims {
		pm := anim.GetParam()
		sortAnims = append(sortAnims, sortParam{
			ID:    id,
			Index: pm.AnimType*100 + pm.PosY*6 + pm.PosX,
		})
	}

	sort.Slice(sortAnims, func(i, j int) bool {
		return sortAnims[i].Index < sortAnims[j].Index
	})

	sortedAnimIDs = []string{}
	for _, a := range sortAnims {
		sortedAnimIDs = append(sortedAnimIDs, a.ID)
	}
}
