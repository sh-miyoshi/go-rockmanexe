package processor

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	swordEndCount = 12
	swordDelay    = 3
)

type Sword struct {
	SkillID      int
	GetObjectPos func(objID string) point.Point
	DamageMgr    *damage.DamageManager
	Arg          skillcore.Argument

	count int
}

func (p *Sword) Process() (bool, error) {
	p.count++

	if p.count == 1*swordDelay {
		// sound.On(resources.SESword) or SEDreamSword // TODO

		dm := damage.Damage{
			DamageType:    damage.TypeObject,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		userPos := p.GetObjectPos(p.Arg.OwnerID)

		targetPos := point.Point{X: userPos.X + 1, Y: userPos.Y}
		if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
			dm.TargetObjID = objID
			p.DamageMgr.New(dm)
		}

		switch p.SkillID {
		case resources.SkillSword:
			// No more damage area
		case resources.SkillWideSword:
			targetPos.Y = userPos.Y - 1
			if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				p.DamageMgr.New(dm)
			}
			targetPos.Y = userPos.Y + 1
			if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				p.DamageMgr.New(dm)
			}
		case resources.SkillLongSword:
			targetPos.X = userPos.X + 2
			if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				p.DamageMgr.New(dm)
			}
		case resources.SkillDreamSword:
			for x := 1; x <= 2; x++ {
				for y := -1; y <= 1; y++ {
					if x == 1 && y == 0 {
						// すでに登録済み
						continue
					}
					if objID := p.Arg.GetPanelInfo(point.Point{X: userPos.X + x, Y: userPos.Y + y}).ObjectID; objID != "" {
						dm.TargetObjID = objID
						p.DamageMgr.New(dm)
					}
				}
			}
		}
	}

	if p.count > swordEndCount {
		return true, nil
	}
	return false, nil
}

func (p *Sword) GetCount() int {
	return p.count
}

func (p *Sword) GetEndCount() int {
	return swordEndCount
}

func (p *Sword) GetSwordType() int {
	switch p.SkillID {
	case resources.SkillSword:
		return 0
	case resources.SkillWideSword:
		return 1
	case resources.SkillLongSword:
		return 2
	case resources.SkillDreamSword:
		return 0
	}
	return 0
}