package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
)

type failed struct {
	ID string
}

func newFailed(objID string) *failed {
	return &failed{
		ID: objID,
	}
}

func (p *failed) Draw() {
}

func (p *failed) Process() (bool, error) {
	// TODO 失敗エフェクトを追加
	return true, nil
}

func (p *failed) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *failed) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
