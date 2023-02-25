package object

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
)

type Player struct {
	objectInfo gameinfo.Object
}

func NewPlayer(info gameinfo.Object) *Player {
	return &Player{
		objectInfo: info,
	}
}

func (p *Player) Process() (bool, error) {
	// TODO
	return false, nil
}

func (p *Player) Draw() {
	// nothing to do at router
}

func (p *Player) DamageProc(dm *damage.Damage) bool {
	// TODO
	return false
}

func (p *Player) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    p.objectInfo.ID,
			Pos:      p.objectInfo.Pos,
			AnimType: anim.AnimTypeObject,
		},
		HP: p.objectInfo.HP,
	}
}

func (p *Player) GetObjectType() int {
	return objanim.ObjTypePlayer
}

func (p *Player) MakeInvisible(count int) {
	// TODO
}
