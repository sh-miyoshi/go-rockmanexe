package deleteanim

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type deleteAction struct {
	id    string
	image int
	pos   point.Point
	count int
}

func New(image int, pos point.Point, isPlayer bool) {
	if isPlayer {
		sound.On(resources.SEPlayerDeleted)
	} else {
		sound.On(resources.SEEnemyDeleted)
	}

	localanim.AnimNew(&deleteAction{
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
	view := battlecommon.ViewPos(p.pos)

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_INVSRC, 255)
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, p.image, true)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
	dxlib.DrawRotaGraph(view.X, view.Y, 1, 0, p.image, true)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
}

func (p *deleteAction) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.id,
		Pos:      p.pos,
		DrawType: anim.DrawTypeEffect,
	}
}
