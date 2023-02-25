package anim

import (
	"fmt"
	"sort"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
)

const (
	AnimTypeObject int = iota + 1
	AnimTypeSkill
	AnimTypeEffect
)

type Param struct {
	ObjID    string
	Pos      common.Point
	AnimType int
}

// Anim ...
type Anim interface {
	Process() (bool, error)
	Draw()
	GetParam() Param
}

var (
	anims         = map[string]Anim{}
	sortedAnimIDs = []string{}
)

func MgrProcess() error {
	for id, anim := range anims {
		end, err := anim.Process()
		if err != nil {
			return fmt.Errorf("Anim process failed: %w", err)
		}

		if end {
			Delete(id)
		}
	}

	sortAnim()

	return nil
}

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
	if _, ok := anims[animID]; !ok {
		return
	}

	delete(anims, animID)
	for i, sid := range sortedAnimIDs {
		if sid == animID {
			sortedAnimIDs = append(sortedAnimIDs[:i], sortedAnimIDs[i+1:]...)
			break
		}
	}
}

func GetAll() []Param {
	res := []Param{}
	for _, anim := range anims {
		res = append(res, anim.GetParam())
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
			Index: pm.AnimType*100 + pm.Pos.Y*6 + pm.Pos.X,
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
