package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type spreadGun struct {
	ID  string
	Arg Argument

	count  int
	drawer skilldraw.DrawSpreadGun
}

type spreadHit struct {
	ID  string
	Arg Argument

	count int
	pos   common.Point
}

func newSpreadGun(objID string, arg Argument) *spreadGun {
	res := &spreadGun{
		ID:  objID,
		Arg: arg,
	}

	return res
}

func (p *spreadGun) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.count)
}

func (p *spreadGun) Process() (bool, error) {
	if p.count == 5 {
		sound.On(resources.SEGun)

		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := common.Point{X: x, Y: pos.Y}
			objs := localanim.ObjAnimGetObjs(objanim.Filter{Pos: &target, ObjType: p.Arg.TargetType})
			if len(objs) > 0 {
				// Hit
				sound.On(resources.SESpreadHit)

				localanim.DamageManager().New(damage.Damage{
					DamageType:    damage.TypeObject,
					Power:         int(p.Arg.Power),
					TargetObjType: p.Arg.TargetType,
					HitEffectType: resources.EffectTypeHitBig,
					Element:       damage.ElementNone,
					TargetObjID:   objs[0].ObjID,
				})
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
							localanim.AnimNew(&spreadHit{
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

	if p.count > resources.SkillSpreadGunEndCount {
		return true, nil
	}
	return false, nil
}

func (p *spreadGun) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *spreadGun) StopByOwner() {
	if p.count < 5 {
		localanim.AnimDelete(p.ID)
	}
}

func (p *spreadHit) Draw() {
}

func (p *spreadHit) Process() (bool, error) {
	p.count++
	if p.count == 10 {
		localanim.AnimNew(effect.Get(resources.EffectTypeSpreadHit, p.pos, 5))
		if objID := field.GetPanelInfo(p.pos).ObjectID; objID != "" {
			localanim.DamageManager().New(damage.Damage{
				DamageType:    damage.TypeObject,
				Power:         int(p.Arg.Power),
				TargetObjType: p.Arg.TargetType,
				HitEffectType: resources.EffectTypeNone,
				Element:       damage.ElementNone,
			})
		}

		return true, nil
	}
	return false, nil
}

func (p *spreadHit) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *spreadHit) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
