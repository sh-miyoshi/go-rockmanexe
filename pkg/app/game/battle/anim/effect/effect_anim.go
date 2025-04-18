package effect

import (
	"sort"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
)

type Anim interface {
	Update() (bool, error)
	Draw()
	GetParam() anim.Param
}

type AnimManager struct {
	anims         map[string]Anim
	sortedAnimIDs []string
}

func NewManager() *AnimManager {
	return &AnimManager{
		anims: make(map[string]Anim),
	}
}

func (am *AnimManager) Update() error {
	for id, anim := range am.anims {
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

func (am *AnimManager) sortAnim() {
	type sortParam struct {
		Index int
		ID    string
	}
	sortAnims := []sortParam{}
	for id, anim := range am.anims {
		pm := anim.GetParam()
		sortAnims = append(sortAnims, sortParam{
			ID:    id,
			Index: pm.Pos.Y*6 + pm.Pos.X,
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
