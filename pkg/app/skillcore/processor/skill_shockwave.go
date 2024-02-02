package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	skilldefines "github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/defines"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type ShockWave struct {
	Arg skillcore.Argument

	count int
	pm    skilldefines.ShockWaveParam
	pos   point.Point
}

func (p *ShockWave) Init(isPlayer bool) {
	p.pm = skilldefines.GetShockWaveParam(isPlayer)
	p.pos = p.Arg.GetObjectPos(p.Arg.OwnerID)
}

func (p *ShockWave) Process() (bool, error) {
	if p.count < p.pm.InitWait {
		p.count++
		return false, nil
	}

	n := p.pm.ImageNum * p.pm.Speed
	if p.count%n == 0 {
		if p.pm.Direct == config.DirectLeft {
			p.pos.X--
		} else if p.pm.Direct == config.DirectRight {
			p.pos.X++
		}

		pn := p.Arg.GetPanelInfo(p.pos)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		p.Arg.DamageMgr.New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           n - 2,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
			Element:       damage.ElementNone,
		})
	}
	p.count++

	if p.pos.X < 0 || p.pos.X > battlecommon.FieldNum.X {
		return true, nil
	}
	return false, nil
}

func (p *ShockWave) GetCount() int {
	return p.count
}

func (p *ShockWave) GetEndCount() int {
	return 6 * 4
}
