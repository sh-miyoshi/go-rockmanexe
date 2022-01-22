package common

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type deleteAction struct {
	id    string
	image int32
	pos   common.Point
	count int
}

func NewDelete(image int32, pos common.Point, isPlayer bool) {
	if isPlayer {
		sound.On(sound.SEPlayerDeleted)
	} else {
		sound.On(sound.SEEnemyDeleted)
	}

	anim.New(&deleteAction{
		id:    uuid.New().String(),
		image: image,
		pos:   pos,
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
	view := ViewPos(p.pos)

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_INVSRC, 255)
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
}

func (p *deleteAction) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.id,
		Pos:      p.pos,
		AnimType: anim.AnimTypeEffect,
	}
}
