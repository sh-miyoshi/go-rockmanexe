package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const (
	shockWaveInitWaitCount = 9
	shockWaveSpeed         = 3
	imgShockWaveNum        = 7
)

type shockWave struct {
	ID    string
	Arg   Argument
	count int
	pos   common.Point
}

func newShockWave(arg Argument) *shockWave {
	return &shockWave{
		ID:  arg.AnimObjID,
		Arg: arg,
		pos: routeranim.ObjAnimGetObjPos(arg.OwnerClientID, arg.OwnerObjectID),
	}
}

func (p *shockWave) Draw() {
	// nothing to do at router
}

func (p *shockWave) Process() (bool, error) {
	p.count++

	if p.count < shockWaveInitWaitCount {
		return false, nil
	}

	n := imgShockWaveNum * shockWaveSpeed
	if p.count%n == 0 {
		p.pos.X++

		pn := p.Arg.GameInfo.GetPanelInfo(p.pos)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		routeranim.DamageNew(p.Arg.OwnerClientID, damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           n - 2,
			TargetType:    p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
			DamageType:    damage.TypeNone,
		})
	}

	if p.pos.X < 0 || p.pos.X > battlecommon.FieldNum.X {
		return true, nil
	}
	return false, nil
}

func (p *shockWave) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeShockWave,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.pos,
		DrawType:  anim.DrawTypeSkill,
		ExtraInfo: info.Marshal(),
	}
}

func (p *shockWave) StopByOwner() {
}

func (p *shockWave) GetEndCount() int {
	return 6 * 4
}
