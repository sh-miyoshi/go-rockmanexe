package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type vulcan struct {
	ID    string
	Arg   Argument
	Times int

	count    int
	atkCount int
	hit      bool
	drawer   skilldraw.DrawVulcan
}

func newVulcan(objID string, arg Argument) *vulcan {
	return &vulcan{
		ID:    objID,
		Arg:   arg,
		Times: 3,
	}
}

func (p *vulcan) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	p.drawer.Draw(view, p.count)
}

func (p *vulcan) Process() (bool, error) {
	p.count++
	if p.count >= resources.SkillVulcanDelay*1 {
		if p.count%(resources.SkillVulcanDelay*5) == resources.SkillVulcanDelay*1 {
			sound.On(resources.SEGun)

			// Add damage
			pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
			hit := false
			p.atkCount++
			lastAtk := p.atkCount == p.Times
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				target := point.Point{X: x, Y: pos.Y}
				if objID := field.GetPanelInfo(target).ObjectID; objID != "" {
					localanim.DamageManager().New(damage.Damage{
						DamageType:    damage.TypeObject,
						Power:         int(p.Arg.Power),
						TargetObjType: p.Arg.TargetType,
						HitEffectType: resources.EffectTypeSpreadHit,
						BigDamage:     lastAtk,
						Element:       damage.ElementNone,
						TargetObjID:   objID,
					})
					localanim.AnimNew(effect.Get(resources.EffectTypeVulcanHit1, target, 20))
					if p.hit && x < battlecommon.FieldNum.X-1 {
						target = point.Point{X: x + 1, Y: pos.Y}
						localanim.AnimNew(effect.Get(resources.EffectTypeVulcanHit2, target, 20))
						if objID := field.GetPanelInfo(target).ObjectID; objID != "" {
							localanim.DamageManager().New(damage.Damage{
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
					sound.On(resources.SECannonHit)
					break
				}
			}
			p.hit = hit
			if lastAtk {
				return true, nil
			}
		}

	}

	return false, nil
}

func (p *vulcan) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *vulcan) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
