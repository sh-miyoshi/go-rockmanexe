package enemy

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
)

type deleteAction struct {
	image int32
	x, y  int
	count int
}

func newDelete(image int32, x, y int) {
	sound.On(sound.SEEnemyDeleted)

	anim.New(&deleteAction{
		image: image,
		x:     x,
		y:     y,
	})
}

func (p *deleteAction) Process() (bool, error) {
	p.count++
	if p.count == 1 {
		anim.New(effect.Get(effect.TypeExplode, p.x, p.y, 0))
	}
	if p.count == 15 {
		dxlib.DeleteGraph(p.image)
		return true, nil
	}
	return false, nil
}

func (p *deleteAction) Draw() {
	x, y := battlecommon.ViewPos(p.x, p.y)

	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_INVSRC, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ADD, 255)
	dxlib.DrawRotaGraph(x, y, 1, 0, p.image, dxlib.TRUE)
	dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 0)
}

func (p *deleteAction) DamageProc(dm *damage.Damage) {
}

func (p *deleteAction) GetParam() anim.Param {
	return anim.Param{
		PosX:     p.x,
		PosY:     p.y,
		AnimType: anim.TypeObject,
		ObjType:  anim.ObjTypeNone,
	}
}
