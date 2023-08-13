package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayHeatShot    = 3
	heatShotAtkDelay = 15
)

const (
	heatShotTypeShot int = iota
	heatShotTypeV
	heatShotTypeSide
)

type heatShot struct {
	ID   string
	Arg  Argument
	Type int

	count int
}

func newHeatShot(objID string, arg Argument, shotType int) *heatShot {
	return &heatShot{
		ID:   objID,
		Arg:  arg,
		Type: shotType,
	}
}

func (p *heatShot) Draw() {
	n := p.count / delayHeatShot

	// Show body
	if n < len(imgHeatShotBody) {
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X+50, view.Y-18, 1, 0, imgHeatShotBody[n], true)
	}

	// Show atk
	n = (p.count - 4) / delayHeatShot
	if n >= 0 && n < len(imgHeatShotAtk) {
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X+100, view.Y-20, 1, 0, imgHeatShotAtk[n], true)
	}
}

func (p *heatShot) Process() (bool, error) {
	if p.count == heatShotAtkDelay {
		sound.On(resources.SEGun)

		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := common.Point{X: x, Y: pos.Y}
			if field.GetPanelInfo(target).ObjectID != "" {
				// TODO: use target object
				// Hit
				localanim.DamageManager().New(damage.Damage{
					DamageType:    damage.TypePosition,
					Pos:           target,
					Power:         int(p.Arg.Power),
					TTL:           1,
					TargetObjType: p.Arg.TargetType,
					HitEffectType: resources.EffectTypeHeatHit,
					Element:       damage.ElementFire,
				})

				// 誘爆
				targets := []common.Point{}
				switch p.Type {
				case heatShotTypeShot:
					targets = append(targets, common.Point{X: target.X + 1, Y: target.Y})
				case heatShotTypeV:
					targets = append(targets, common.Point{X: target.X + 1, Y: target.Y - 1})
					targets = append(targets, common.Point{X: target.X + 1, Y: target.Y + 1})
				case heatShotTypeSide:
					targets = append(targets, common.Point{X: target.X, Y: target.Y - 1})
					targets = append(targets, common.Point{X: target.X, Y: target.Y + 1})
				}

				for _, t := range targets {
					localanim.AnimNew(effect.Get(resources.EffectTypeHeatHit, t, 0))
					// TODO: use target object
					localanim.DamageManager().New(damage.Damage{
						DamageType:    damage.TypePosition,
						Pos:           t,
						Power:         int(p.Arg.Power),
						TTL:           1,
						TargetObjType: p.Arg.TargetType,
						HitEffectType: resources.EffectTypeNone,
						Element:       damage.ElementFire,
					})
				}

				break
			}
		}
	}

	p.count++

	max := len(imgHeatShotAtk)
	if len(imgHeatShotBody) > max {
		max = len(imgHeatShotBody)
	}

	if p.count > max*delayHeatShot {
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
