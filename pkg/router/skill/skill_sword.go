package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const (
	TypeSword int = iota
	TypeWideSword
	TypeLongSword

	TypeSwordMax
)

const (
	delaySword = 3
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

	if p.count == 1*delaySword {
		dm := damage.Damage{
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: 0, // TODO: 正しい値をセット
			BigDamage:     true,
			DamageType:    damage.TypeNone,
		}

		pos := objanim.GetObjPos(p.Arg.OwnerObjectID)

		dm.Pos.X = pos.X + 1
		dm.Pos.Y = pos.Y
		damage.New(dm)

		switch p.Type {
		case TypeSword:
			// No more damage area
		case TypeWideSword:
			dm.Pos.Y = pos.Y - 1
			damage.New(dm)
			dm.Pos.Y = pos.Y + 1
			damage.New(dm)
		case TypeLongSword:
			dm.Pos.X = pos.X + 2
			damage.New(dm)
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
	case TypeSword:
		info.AnimType = routeranim.TypeSword
	case TypeWideSword:
		info.AnimType = routeranim.TypeWideSword
	case TypeLongSword:
		info.AnimType = routeranim.TypeLongSword
	}

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       objanim.GetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *sword) StopByOwner() {
	anim.Delete(p.ID)
}

func (p *sword) GetEndCount() int {
	const imgSwordNum = 4

	return imgSwordNum * delaySword
}