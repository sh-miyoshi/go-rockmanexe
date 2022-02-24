package objanim

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/stretchr/stew/slice"
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
	MakeInvisible(count int)
}

var (
	anims         = map[string]Anim{}
	sortedAnimIDs = []string{}
	activeAnimIDs = []string{}
)

// MgrProcess ...
func MgrProcess(enableDamage, blackout bool) error {
	for id, anim := range anims {
		if blackout && !slice.Contains(activeAnimIDs, id) {
			continue
		}

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
			if dm := damage.Get(pm.Pos); dm != nil {
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
	id := anim.GetParam().ObjID
	if id == "" {
		id = uuid.New().String()
	}

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
	activeAnimIDs = []string{}
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

func GetObjPos(objID string) common.Point {
	for _, anim := range anims {
		pm := anim.GetParam()
		if pm.ObjID == objID {
			return pm.Pos
		}
	}

	return common.Point{X: -1, Y: -1}
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

func AddActiveAnim(id string) {
	activeAnimIDs = append(activeAnimIDs, id)
}

func MakeInvisible(id string, count int) {
	logger.Debug("ID: %s, count: %d, anims: %+v", id, count, anims)
	if _, ok := anims[id]; ok {
		anims[id].MakeInvisible(count)
	}
}

func ExistsObject(pos common.Point) string {
	for _, anim := range anims {
		pm := anim.GetParam()
		if pm.Pos == pos {
			return pm.ObjID
		}
	}
	return ""
}

func sortAnim() {
	type sortParam struct {
		Index int
		ID    string
	}
	sortAnims := []sortParam{}
	for id, anim := range anims {
		pm := anim.GetParam()
		index := pm.Pos.Y*6 + pm.Pos.X
		if slice.Contains(activeAnimIDs, id) {
			index += 100
		}

		sortAnims = append(sortAnims, sortParam{
			ID:    id,
			Index: index,
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
