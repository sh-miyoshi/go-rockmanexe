package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type ShockWaveParam struct {
	InitWait int
	Speed    int
	Direct   int
	ImageNum int
}

type ShockWave struct {
	Arg skillcore.Argument

	count int
	pm    ShockWaveParam
	pos   point.Point
}

func (p *ShockWave) Init(isPlayer bool) {
	if isPlayer {
		p.pm = ShockWaveParam{
			InitWait: 9,
			Speed:    3,
			Direct:   config.DirectRight,
			ImageNum: 9,
		}
	} else {
		p.pm = ShockWaveParam{
			InitWait: 0,
			Speed:    5,
			Direct:   config.DirectLeft,
			ImageNum: 9,
		}
	}
	p.pos = p.Arg.GetObjectPos(p.Arg.OwnerID)
}

func (p *ShockWave) Process() (bool, error) {
	if p.count < p.pm.InitWait {
		p.count++
		return false, nil
	}

	n := p.pm.ImageNum * p.pm.Speed
	if p.count%n == 0 {
		p.Arg.SoundOn(resources.SEShockWave)
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

func (p *ShockWave) GetParam() ShockWaveParam {
	return p.pm
}
