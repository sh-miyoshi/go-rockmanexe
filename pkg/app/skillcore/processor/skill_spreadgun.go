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
	DamageMgr *damage.DamageManager
	Arg       skillcore.Argument
	Pos       point.Point
	count     int
}

type SpreadGun struct {
	GetObjectPos func(objID string) point.Point
	DamageMgr    *damage.DamageManager
	Arg          skillcore.Argument

	count int
	hits  []SpreadHit
}

func (p *SpreadGun) Process() (bool, error) {
	if p.count == 5 {
		// sound.On(resources.SEGun) // TODO

		pos := p.GetObjectPos(p.Arg.OwnerID)
		dm := damage.Damage{
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
				// sound.On(resources.SESpreadHit) // TODO

				dm.TargetObjID = objID
				logger.Debug("Add damage by spread gun: %+v", dm)
				p.DamageMgr.New(dm)

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
								DamageMgr: p.DamageMgr,
								Arg:       p.Arg,
								Pos:       pos,
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

func (p *SpreadGun) GetEndCount() int {
	return spreadGunEndCount
}

func (p *SpreadGun) PopSpreadHits() []SpreadHit {
	if len(p.hits) > 0 {
		res := append([]SpreadHit{}, p.hits...)
		p.hits = []SpreadHit{}
		return res
	}
	return []SpreadHit{}
}

func (p *SpreadHit) Process() (bool, error) {
	if p.count == 10 {
		if objID := p.Arg.GetPanelInfo(p.Pos).ObjectID; objID != "" {
			p.DamageMgr.New(damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				Element:       damage.ElementNone,
				TargetObjID:   objID,
			})
		}

		return true, nil
	}
	p.count++
	return false, nil
}

func (p *SpreadHit) GetCount() int {
	return p.count
}