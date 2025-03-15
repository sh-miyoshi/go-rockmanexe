package skill

import (
	"sort"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/stretchr/stew/slice"
)

type Anim interface {
	Update() (bool, error)
	Draw()
	GetParam() anim.Param
}

type AnimManager struct {
	anims         map[string]Anim
	sortedAnimIDs []string
	activeAnimIDs []string
}

func NewManager() *AnimManager {
	return &AnimManager{
		anims: make(map[string]Anim),
	}
}

func (am *AnimManager) Update(isActive bool) error {
	for id, anim := range am.anims {
		if !isActive && !slice.Contains(am.activeAnimIDs, id) {
			continue
		}

		end, err := anim.Update()
		if err != nil {
			return errors.Wrap(err, "Anim process failed")
		}

		if end {
			am.Delete(id)
		}
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
	id := uuid.New().String()
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
}

func (am *AnimManager) Delete(animID string) {
	if _, ok := am.anims[animID]; !ok {
		return
	}

	delete(am.anims, animID)
	for i, sid := range am.sortedAnimIDs {
		if sid == animID {
			am.sortedAnimIDs = append(am.sortedAnimIDs[:i], am.sortedAnimIDs[i+1:]...)
			break
		}
	}
}

func (am *AnimManager) GetAll() []anim.Param {
	res := []anim.Param{}
	for _, anim := range am.anims {
		res = append(res, anim.GetParam())
	}

	return res
}

func (am *AnimManager) SetActiveAnim(id string) {
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
