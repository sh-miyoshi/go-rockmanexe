package objanim

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
)

const (
	ObjTypePlayer int = 1 << iota
	ObjTypeEnemy
	ObjTypeNone
)

const (
	ObjTypeAll = ObjTypePlayer | ObjTypeEnemy | ObjTypeNone
)

type Filter struct {
	ObjID   string
	ObjType int
}

// Anim ...
type Anim interface {
	Process() (bool, error)
	Draw()
	DamageProc(dm *damage.Damage) bool
	GetParam() anim.Param
	GetObjectType() int
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
			Delete(id)
		}
	}

	// Damage Process
	if enableDamage {
		hit := []string{}
		for _, anim := range anims {
			pm := anim.GetParam()
			if dm := damage.Get(pm.PosX, pm.PosY); dm != nil {
				if anim.DamageProc(dm) {
					hit = append(hit, dm.ID)
				}
			}
		}

		for _, h := range hit {
			damage.Remove(h)
		}

		damage.MgrProcess()
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
func IsProcessing(animID string) bool {
	_, exists := anims[animID]
	return exists
}

func Cleanup() {
	anims = map[string]Anim{}
	sortedAnimIDs = []string{}
}

func Delete(animID string) {
	delete(anims, animID)
	for i, sid := range sortedAnimIDs {
		if sid == animID {
			sortedAnimIDs = append(sortedAnimIDs[:i], sortedAnimIDs[i+1:]...)
			break
		}
	}
}

func GetObjPos(objID string) (x, y int) {
	for _, anim := range anims {
		pm := anim.GetParam()
		if pm.ObjID == objID {
			return pm.PosX, pm.PosY
		}
	}

	return -1, -1
}

func GetObjs(filter Filter) []anim.Param {
	res := []anim.Param{}
	if filter.ObjID != "" {
		for _, anim := range anims {
			pm := anim.GetParam()
			res = append(res, pm)
		}
		return res
	}

	for _, anim := range anims {
		pm := anim.GetParam()
		if filter.ObjType&anim.GetObjectType() != 0 {
			res = append(res, pm)
		}
	}

	return res
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
			Index: pm.PosY*6 + pm.PosX,
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
