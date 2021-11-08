package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type miniBomb struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int
	TargetX    int
	TargetY    int

	count int
	dist  int
	baseX int32
	baseY int32
	dx    int
	dy    int
}

func (p *miniBomb) Draw() {
	n := (p.count / delayMiniBomb) % len(imgMiniBomb)
	if n >= 0 {
		vx := p.baseX + int32(p.dx)
		vy := p.baseY + int32(p.dy)
		dxlib.DrawRotaGraph(vx-38, vy, 1, 0, imgMiniBomb[n], dxlib.TRUE)
	}
}

func (p *miniBomb) Process() (bool, error) {
	if p.count == 0 {
		// Initialize
		px, py := objanim.GetObjPos(p.OwnerID)
		p.baseX, p.baseY = battlecommon.ViewPos(px, py)
		// TODO: yが等しい場合でかつプレイヤー側のみ
		p.dist = (p.TargetX - px) * field.PanelSizeX

		sound.On(sound.SEBombThrow)
	}

	// y = ax^2 + bx +c
	// (0,0), (d/2, ymax), (d, 0)
	p.count++
	p.dx += 4
	ymax := 100
	p.dy = ymax*4*p.dx*p.dx/(p.dist*p.dist) - ymax*4*p.dx/p.dist

	if p.dx >= p.dist+38 {
		// TODO 不発処理(画面外やパネル状況など)
		sound.On(sound.SEExplode)
		anim.New(effect.Get(effect.TypeExplode, p.TargetX, p.TargetY, 0))
		damage.New(damage.Damage{
			PosX:          p.TargetX,
			PosY:          p.TargetY,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			BigDamage:     true,
		})
		return true, nil
	}
	return false, nil
}

func (p *miniBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
