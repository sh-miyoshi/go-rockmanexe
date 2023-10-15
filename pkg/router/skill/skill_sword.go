package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type sword struct {
	ID   string
	Type int
	Arg  Argument

	count int
}

func newSword(swordType int, arg Argument) *sword {
	return &sword{
		ID:   arg.AnimObjID,
		Type: swordType,
		Arg:  arg,
	}
}

func (p *sword) Draw() {
	// nothing to do at router
}

func (p *sword) Process() (bool, error) {
	p.count++

	if p.count == 1*resources.SkillSwordDelay {
		dm := damage.Damage{
			DamageType:    damage.TypeObject,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		userPos := routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID)

		targetPos := common.Point{X: userPos.X + 1, Y: userPos.Y}
		if objID := p.Arg.GameInfo.GetPanelInfo(targetPos).ObjectID; objID != "" {
			dm.TargetObjID = objID
			routeranim.DamageNew(p.Arg.OwnerClientID, dm)
		}

		switch p.Type {
		case resources.SkillTypeSword:
			// No more damage area
		case resources.SkillTypeWideSword:
			targetPos.Y = userPos.Y - 1
			if objID := p.Arg.GameInfo.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				routeranim.DamageNew(p.Arg.OwnerClientID, dm)
			}
			targetPos.Y = userPos.Y + 1
			if objID := p.Arg.GameInfo.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
			}
		case resources.SkillTypeLongSword:
			targetPos.X = userPos.X + 2
			if objID := p.Arg.GameInfo.GetPanelInfo(targetPos).ObjectID; objID != "" {
				dm.TargetObjID = objID
				routeranim.DamageNew(p.Arg.OwnerClientID, dm)
			}
		}
	}

	if p.count > p.GetEndCount() {
		return true, nil
	}
	return false, nil
}

func (p *sword) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.count,
	}
	switch p.Type {
	case resources.SkillTypeSword:
		info.AnimType = routeranim.TypeSword
	case resources.SkillTypeWideSword:
		info.AnimType = routeranim.TypeWideSword
	case resources.SkillTypeLongSword:
		info.AnimType = routeranim.TypeLongSword
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *sword) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *sword) GetEndCount() int {
	return resources.SkillSwordEndCount
}
