package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

const ()

type spreadGun struct {
	ID    string
	Arg   Argument
	count int
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
	}
}

func (p *spreadGun) Draw() {
	// nothing to do at router
}

func (p *spreadGun) Process() (bool, error) {
	if p.count == 5 {
		pos := routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID)
		dm := damage.Damage{
			OwnerClientID: p.Arg.OwnerClientID,
			Pos:           pos,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: battlecommon.EffectTypeHitBig,
			BigDamage:     true,
			DamageType:    damage.TypeNone,
		}

		if p.Arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				dm.Pos.X = x
				if p.Arg.GameInfo.GetPanelInfo(common.Point{X: x, Y: dm.Pos.Y}).ObjectID != "" {
					logger.Debug("Add damage by spread gun: %+v", dm)
					routeranim.DamageNew(p.Arg.OwnerClientID, dm)

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
								routeranim.AnimNew(p.Arg.OwnerClientID, &spreadHit{
									Arg: p.Arg,
									pos: common.Point{X: x + sx, Y: pos.Y + sy},
								})
							}
						}
					}
					break
				}
			}
		} else {
			logger.Error("unexpected target type for spread gun: %+v", p.Arg)
			return false, fmt.Errorf("unexpected target type")
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
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		DrawType:  anim.DrawTypeEffect,
		ExtraInfo: info.Marshal(),
	}
}

func (p *spreadGun) StopByOwner() {
	if p.count < 5 {
		routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
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
		routeranim.DamageNew(p.Arg.OwnerClientID, damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: battlecommon.EffectTypeNone,
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
