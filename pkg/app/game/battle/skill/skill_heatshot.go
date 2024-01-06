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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	heatShotAtkDelay = 15
)

const (
	heatShotTypeShot int = iota
	heatShotTypeV
	heatShotTypeSide
)

type heatShot struct {
	ID   string
	Arg  skillcore.Argument
	Type int

	count  int
	drawer skilldraw.DrawHeatShot
}

func newHeatShot(objID string, arg skillcore.Argument, shotType int) *heatShot {
	return &heatShot{
		ID:   objID,
		Arg:  arg,
		Type: shotType,
	}
}

func (p *heatShot) Draw() {
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)
	p.drawer.Draw(view, p.count)
}

func (p *heatShot) Process() (bool, error) {
	if p.count == heatShotAtkDelay {
		sound.On(resources.SEGun)

		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := point.Point{X: x, Y: pos.Y}
			if objID := field.GetPanelInfo(target).ObjectID; objID != "" {
				// Hit
				localanim.DamageManager().New(damage.Damage{
					DamageType:    damage.TypeObject,
					Power:         int(p.Arg.Power),
					TargetObjType: p.Arg.TargetType,
					HitEffectType: resources.EffectTypeHeatHit,
					Element:       damage.ElementFire,
					TargetObjID:   objID,
				})

				// 誘爆
				targets := []point.Point{}
				switch p.Type {
				case heatShotTypeShot:
					targets = append(targets, point.Point{X: target.X + 1, Y: target.Y})
				case heatShotTypeV:
					targets = append(targets, point.Point{X: target.X + 1, Y: target.Y - 1})
					targets = append(targets, point.Point{X: target.X + 1, Y: target.Y + 1})
				case heatShotTypeSide:
					targets = append(targets, point.Point{X: target.X, Y: target.Y - 1})
					targets = append(targets, point.Point{X: target.X, Y: target.Y + 1})
				}

				for _, t := range targets {
					localanim.AnimNew(effect.Get(resources.EffectTypeHeatHit, t, 0))
					if objID := field.GetPanelInfo(t).ObjectID; objID != "" {
						localanim.DamageManager().New(damage.Damage{
							DamageType:    damage.TypeObject,
							Power:         int(p.Arg.Power),
							TargetObjType: p.Arg.TargetType,
							HitEffectType: resources.EffectTypeNone,
							Element:       damage.ElementFire,
							TargetObjID:   objID,
						})
					}
				}

				break
			}
		}
	}

	p.count++

	if p.count > resources.SkillHeatShotEndCount {
		return true, nil
	}
	return false, nil
}

func (p *heatShot) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *heatShot) StopByOwner() {
	if p.count < heatShotAtkDelay {
		localanim.AnimDelete(p.ID)
	}
}
