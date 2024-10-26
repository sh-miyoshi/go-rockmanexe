package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	spreadGunEndCount = 8 // imgAtkNum*delay
)

type SpreadHit struct {
	Arg   skillcore.Argument
	Pos   point.Point
	count int
}

type SpreadGun struct {
	Arg skillcore.Argument

	count int
	hits  []SpreadHit
}

func (p *SpreadGun) Update() (bool, error) {
	if p.count == 5 {
		p.Arg.SoundOn(resources.SEGun)

		pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
		dm := damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypeObject,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeHitBig,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := point.Point{X: x, Y: pos.Y}
			if objID := p.Arg.GetPanelInfo(target).ObjectID; objID != "" {
				// Hit
				p.Arg.SoundOn(resources.SESpreadHit)

				dm.TargetObjID = objID
				logger.Debug("Add damage by spread gun: %+v", dm)
				p.Arg.DamageMgr.New(dm)

				// Spreading
				for sy := -1; sy <= 1; sy++ {
					if pos.Y+sy < 0 || pos.Y+sy >= battlecommon.FieldNum.Y {
						continue
					}
					for sx := -1; sx <= 1; sx++ {
						if sy == 0 && sx == 0 {
							continue
						}
						if x+sx >= 0 && x+sx < battlecommon.FieldNum.X {
							pos := point.Point{X: x + sx, Y: pos.Y + sy}
							logger.Debug("Add spread hit to %s", pos.String())
							p.hits = append(p.hits, SpreadHit{
								Arg: p.Arg,
								Pos: pos,
							})
						}
					}
				}

				break
			}
		}
	}

	p.count++

	if p.count > spreadGunEndCount {
		return true, nil
	}
	return false, nil
}

func (p *SpreadGun) GetCount() int {
	return p.count
}

func (p *SpreadGun) PopSpreadHits() []SpreadHit {
	if len(p.hits) > 0 {
		res := append([]SpreadHit{}, p.hits...)
		p.hits = []SpreadHit{}
		return res
	}
	return []SpreadHit{}
}

func (p *SpreadHit) Update() (bool, error) {
	p.count++
	if p.count == 1 {
		if objID := p.Arg.GetPanelInfo(p.Pos).ObjectID; objID != "" {
			p.Arg.DamageMgr.New(damage.Damage{
				OwnerClientID: p.Arg.OwnerClientID,
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				Element:       damage.ElementNone,
				TargetObjID:   objID,
			})
		}
	}

	return p.count >= 10, nil
}

func (p *SpreadHit) GetCount() int {
	return p.count
}
