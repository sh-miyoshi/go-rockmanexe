package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

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
			DamageType:    damage.TypeObject,
			OwnerClientID: p.Arg.OwnerClientID,
			Power:         int(p.Arg.Power),
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeHitBig,
			BigDamage:     true,
			Element:       damage.ElementNone,
		}

		if p.Arg.TargetType == damage.TargetEnemy {
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				if objID := p.Arg.GameInfo.GetPanelInfo(common.Point{X: x, Y: pos.Y}).ObjectID; objID != "" {
					dm.TargetObjID = objID
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
	return resources.SkillSpreadGunEndCount
}

func (p *spreadHit) Draw() {
	// nothing to do at router
}

func (p *spreadHit) Process() (bool, error) {
	p.count++
	if p.count == 10 {
		if objID := p.Arg.GameInfo.GetPanelInfo(p.pos).ObjectID; objID != "" {
			routeranim.DamageNew(p.Arg.OwnerClientID, damage.Damage{
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
