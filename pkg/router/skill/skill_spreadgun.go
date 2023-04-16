package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const ()

type spreadGun struct {
	ID    string
	Arg   Argument
	count int
	pos   common.Point
}

type spreadHit struct {
	ID  string
	Arg Argument

	count int
	pos   common.Point
}

func newSpreadGun(arg Argument) *spreadGun {
	return &spreadGun{
		ID:  arg.AnimObjID,
		Arg: arg,
		pos: common.Point{X: -1, Y: -1},
	}
}

func (p *spreadGun) Draw() {
	// nothing to do at router
}

func (p *spreadGun) Process() (bool, error) {
	if p.count == 5 {
		pos := objanim.GetObjPos(p.Arg.OwnerObjectID)
		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := common.Point{X: x, Y: pos.Y}
			objs := objanim.GetObjs(objanim.Filter{Pos: &target, ObjType: p.Arg.TargetType})
			if len(objs) > 0 {
				// Hit
				damage.New(damage.Damage{
					Pos:           target,
					Power:         int(p.Arg.Power),
					TTL:           1,
					TargetType:    p.Arg.TargetType,
					HitEffectType: 0, // TODO: 正しい値をセット
					DamageType:    damage.TypeNone,
				})
				p.pos = target

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
							anim.New(&spreadHit{
								Arg: p.Arg,
								pos: common.Point{X: x + sx, Y: pos.Y + sy},
							})
						}
					}
				}

				break
			}
		}
	}

	p.count++

	if p.count > p.GetEndCount() {
		return true, nil
	}
	return false, nil
}

func (p *spreadGun) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeSpreadGun,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.pos,
		DrawType:  anim.DrawTypeEffect,
		ExtraInfo: info.Marshal(),
	}
}

func (p *spreadGun) StopByOwner() {
	if p.count < 5 {
		anim.Delete(p.ID)
	}
}

func (p *spreadGun) GetEndCount() int {
	const (
		delaySpreadGun      = 2
		imgSpreadGunBodyNum = 4
		imgSpreadGunAtkNum  = 4
	)

	max := imgSpreadGunAtkNum
	if imgSpreadGunBodyNum > max {
		max = imgSpreadGunBodyNum
	}

	return max * delaySpreadGun
}

func (p *spreadHit) Draw() {
	// nothing to do at router
}

func (p *spreadHit) Process() (bool, error) {
	p.count++
	if p.count == 10 {
		// TODO: add effect anim
		damage.New(damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: 0, // TODO: 正しい値をセット
			DamageType:    damage.TypeNone,
		})

		return true, nil
	}
	return false, nil
}

func (p *spreadHit) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeSpreadHit,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.count,
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.pos,
		DrawType:  anim.DrawTypeEffect,
		ExtraInfo: info.Marshal(),
	}
}
