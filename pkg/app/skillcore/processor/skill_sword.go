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
	SkillID int
	Arg     skillcore.Argument

	count int
}

func (p *Sword) Process() (bool, error) {
	p.count++

	if p.count == 1*swordDelay {
		if p.SkillID == resources.SkillDreamSword {
			p.Arg.SoundOn(resources.SEDreamSword)
		} else {
			p.Arg.SoundOn(resources.SESword)
		}

		dm := damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			DamageType:    damage.TypeObject,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		userPos := p.Arg.GetObjectPos(p.Arg.OwnerID)

		targetPos := point.Point{X: userPos.X + 1, Y: userPos.Y}
		if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
			dm.TargetObjID = objID
			p.Arg.DamageMgr.New(dm)
		}

		switch p.SkillID {
		case resources.SkillSword:
			// No more damage area
		case resources.SkillWideSword:
			targetPos.Y = userPos.Y - 1
			if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				p.Arg.DamageMgr.New(dm)
			}
			targetPos.Y = userPos.Y + 1
			if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				p.Arg.DamageMgr.New(dm)
			}
		case resources.SkillLongSword:
			targetPos.X = userPos.X + 2
			if objID := p.Arg.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				p.Arg.DamageMgr.New(dm)
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
						p.Arg.DamageMgr.New(dm)
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

func (p *Sword) GetDelay() int {
	return swordDelay
}
