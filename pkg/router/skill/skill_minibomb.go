package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const (
	miniBombEndCount   = 60
	delayMiniBombThrow = 4
)

type miniBomb struct {
	ID  string
	Arg Argument

	count  int
	pos    common.Point
	target common.Point
}

func newMiniBomb(arg Argument) *miniBomb {
	pos := objanim.GetObjPos(arg.OwnerID)
	return &miniBomb{
		ID:     arg.AnimObjID,
		Arg:    arg,
		pos:    pos,
		target: common.Point{X: pos.X + 3, Y: pos.Y},
	}
}

func (p *miniBomb) Draw() {
	// nothing to do at router
}

func (p *miniBomb) Process() (bool, error) {
	p.count++

	if p.count == miniBombEndCount {
		pn := p.Arg.GameInfo.GetPanelInfo(p.target)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		// TODO: effect
		damage.New(damage.Damage{
			OwnerClientID: p.Arg.OwnerID,
			Pos:           p.target,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: 0, // TODO: 正しい値をセット
			BigDamage:     true,
			DamageType:    damage.TypeNone,
		})
		return true, nil
	}
	return false, nil
}

func (p *miniBomb) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeMiniBomb,
		OwnerClientID: p.Arg.OwnerID,
		ActCount:      p.count,
	}

	// TODO: pos

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		ExtraInfo: info.Marshal(),
	}
}

func (p *miniBomb) StopByOwner() {
	if p.count < 5 {
		anim.Delete(p.ID)
	}
}

func (p *miniBomb) GetEndCount() int {
	return miniBombEndCount
}