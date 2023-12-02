package objanim

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
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
	Pos     *point.Point
}

type Param struct {
	anim.Param

	HP int
}

type Anim interface {
	Process() (bool, error)
	Draw()
	DamageProc(dm *damage.Damage) bool
	GetParam() Param
	GetObjectType() int
	MakeInvisible(count int)
}

type AnimManager struct {
	anims         map[string]Anim
	sortedAnimIDs []string
	activeAnimIDs []string
	dmMgr         *damage.DamageManager
}

var (
	FilterAll = Filter{ObjType: ObjTypeAll}
)

func NewManager() *AnimManager {
	return &AnimManager{
		anims: make(map[string]Anim),
		dmMgr: damage.NewManager(),
	}
}

func (am *AnimManager) Process(enableDamage, blackout bool) error {
	for id, anim := range am.anims {
		if blackout && !slice.Contains(am.activeAnimIDs, id) {
			continue
		}

		end, err := anim.Process()
		if err != nil {
			return fmt.Errorf("Anim process failed: %w", err)
		}

		if end {
			am.Delete(id)
		}
	}

	// Damage Process
	if enableDamage {
		hit := []string{}
		for _, anim := range am.anims {
			pm := anim.GetParam()
			if dm := am.dmMgr.GetHitDamage(pm.Pos, pm.ObjID); dm != nil {
				if anim.DamageProc(dm) {
					hit = append(hit, dm.ID)
				}
			}
		}

		if len(hit) > 0 {
			logger.Debug("Hit damages: %+v", hit)
			for _, h := range hit {
				am.dmMgr.Remove(h)
			}
		}

		am.dmMgr.Process()
	}

	am.sortAnim()

	return nil
}

func (am *AnimManager) Draw() {
	for _, id := range am.sortedAnimIDs {
		am.anims[id].Draw()
	}
}

func (am *AnimManager) New(anim Anim) string {
	id := anim.GetParam().ObjID
	if id == "" {
		id = uuid.New().String()
	}

	am.anims[id] = anim
	am.sortAnim()
	return id
}

func (am *AnimManager) IsProcessing(animID string) bool {
	_, exists := am.anims[animID]
	return exists
}

func (am *AnimManager) Cleanup() {
	am.anims = map[string]Anim{}
	am.sortedAnimIDs = []string{}
	am.activeAnimIDs = []string{}
	am.dmMgr.RemoveAll()
}

func (am *AnimManager) Delete(animID string) {
	delete(am.anims, animID)
	for i, sid := range am.sortedAnimIDs {
		if sid == animID {
			am.sortedAnimIDs = append(am.sortedAnimIDs[:i], am.sortedAnimIDs[i+1:]...)
			break
		}
	}
}

func (am *AnimManager) GetObjPos(objID string) point.Point {
	for _, anim := range am.anims {
		pm := anim.GetParam()
		if pm.ObjID == objID {
			return pm.Pos
		}
	}

	return point.Point{X: -1, Y: -1}
}

func (am *AnimManager) GetObjs(filter Filter) []Param {
	res := []Param{}

	for _, anim := range am.anims {
		pm := anim.GetParam()
		if filter.ObjID != "" && pm.ObjID != filter.ObjID {
			continue
		}
		if filter.Pos != nil && (pm.Pos.X != filter.Pos.X || pm.Pos.Y != filter.Pos.Y) {
			continue
		}
		if filter.ObjType&anim.GetObjectType() == 0 {
			continue
		}
		res = append(res, pm)
	}

	return res
}

func (am *AnimManager) AddActiveAnim(id string) {
	am.activeAnimIDs = append(am.activeAnimIDs, id)
}

func (am *AnimManager) DeactivateAnim(id string) {
	animIDs := []string{}
	for _, animID := range am.activeAnimIDs {
		if id != animID {
			animIDs = append(animIDs, animID)
		}
	}
	am.activeAnimIDs = animIDs
}

func (am *AnimManager) MakeInvisible(id string, count int) {
	logger.Debug("ID: %s, count: %d, anims: %+v", id, count, am.anims)
	if _, ok := am.anims[id]; ok {
		am.anims[id].MakeInvisible(count)
	}
}

func (am *AnimManager) ExistsObject(pos point.Point) string {
	objs := am.GetObjs(Filter{Pos: &pos, ObjType: ObjTypeAll})
	if len(objs) > 0 {
		return objs[0].ObjID
	}

	return ""
}

func (am *AnimManager) DamageManager() *damage.DamageManager {
	return am.dmMgr
}

func (am *AnimManager) sortAnim() {
	type sortParam struct {
		Index int
		ID    string
	}
	sortAnims := []sortParam{}
	for id, anim := range am.anims {
		pm := anim.GetParam()
		index := pm.Pos.Y*6 + pm.Pos.X
		if slice.Contains(am.activeAnimIDs, id) {
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

	am.sortedAnimIDs = []string{}
	for _, a := range sortAnims {
		am.sortedAnimIDs = append(am.sortedAnimIDs, a.ID)
	}
}
