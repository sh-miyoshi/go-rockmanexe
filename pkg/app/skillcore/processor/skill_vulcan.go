package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	vulcanDelay = 2
)

type Vulcan struct {
	Arg   skillcore.Argument
	Times int

	count    int
	isHit    bool
	atkCount int
	effects  []resources.EffectParam
}

func (p *Vulcan) Update() (bool, error) {
	p.count++
	if p.count >= vulcanDelay*1 {
		if p.count%(vulcanDelay*5) == vulcanDelay*1 {
			p.Arg.SoundOn(resources.SEGun)

			// Add damage
			pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
			hit := false
			p.atkCount++
			lastAtk := p.atkCount == p.Times
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				target := point.Point{X: x, Y: pos.Y}
				if objID := p.Arg.GetPanelInfo(target).ObjectID; objID != "" {
					p.Arg.DamageMgr.New(damage.Damage{
						OwnerClientID: p.Arg.OwnerClientID,
						DamageType:    damage.TypeObject,
						Power:         int(p.Arg.Power),
						TargetObjType: p.Arg.TargetType,
						HitEffectType: resources.EffectTypeSpreadHit,
						BigDamage:     lastAtk,
						Element:       damage.ElementNone,
						TargetObjID:   objID,
					})
					p.effects = append(p.effects, resources.EffectParam{
						Type:      resources.EffectTypeVulcanHit1,
						Pos:       target,
						RandRange: 20,
					})
					if p.isHit && x < battlecommon.FieldNum.X-1 {
						target = point.Point{X: x + 1, Y: pos.Y}
						p.effects = append(p.effects, resources.EffectParam{
							Type:      resources.EffectTypeVulcanHit2,
							Pos:       target,
							RandRange: 20,
						})
						if objID := p.Arg.GetPanelInfo(target).ObjectID; objID != "" {
							p.Arg.DamageMgr.New(damage.Damage{
								OwnerClientID: p.Arg.OwnerClientID,
								DamageType:    damage.TypeObject,
								Power:         int(p.Arg.Power),
								TargetObjType: p.Arg.TargetType,
								HitEffectType: resources.EffectTypeNone,
								BigDamage:     lastAtk,
								Element:       damage.ElementNone,
								TargetObjID:   objID,
							})
						}
					}
					hit = true
					p.Arg.SoundOn(resources.SECannonHit)
					break
				}
			}
			p.isHit = hit
			if lastAtk {
				return true, nil
			}
		}

	}

	return false, nil
}

func (p *Vulcan) GetCount() int {
	return p.count
}

func (p *Vulcan) PopEffects() []resources.EffectParam {
	if len(p.effects) > 0 {
		res := append([]resources.EffectParam{}, p.effects...)
		p.effects = []resources.EffectParam{}
		return res
	}
	return []resources.EffectParam{}
}

func (p *Vulcan) GetDelay() int {
	return vulcanDelay
}
