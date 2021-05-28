package common

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type deleteAction struct {
	id    string
	image int32
	x, y  int
	count int
}

func NewDelete(image int32, x, y int, isPlayer bool) {
	if isPlayer {
		sound.On(sound.SEPlayerDeleted)
	} else {
		sound.On(sound.SEEnemyDeleted)
	}

	anim.New(&deleteAction{
		id:    uuid.New().String(),
		image: image,
		x:     x,
		y:     y,
	})
}

func (p *deleteAction) Process() (bool, error) {
	p.count++
	if p.count == 15 {
		dxlib.DeleteGraph(p.image)
		return true, nil
	}
	return false, nil
}

func (p *deleteAction) Draw() {
	x, y := ViewPos(p.x, p.y)

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_INVSRC, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
}

func (p *deleteAction) DamageProc(dm *damage.Damage) bool {
	return false
}

func (p *deleteAction) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.id,
		PosX:     p.x,
		PosY:     p.y,
		AnimType: anim.TypeEffect,
		ObjType:  anim.ObjTypeNone,
	}
}
