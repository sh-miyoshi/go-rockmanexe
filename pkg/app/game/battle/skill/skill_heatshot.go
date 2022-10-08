package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayHeatShot = 2
)

type heatShot struct {
	ID  string
	Arg Argument

	count int
}

func newHeatShot(objID string, arg Argument) *heatShot {
	return &heatShot{
		ID:  objID,
		Arg: arg,
	}
}

func (p *heatShot) Draw() {
	n := p.count / delayHeatShot

	// Show body
	if n < len(imgHeatShotBody) {
		pos := objanim.GetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X+50, view.Y-18, 1, 0, imgHeatShotBody[n], true)
	}

	// Show atk
	n = (p.count - 4) / delayHeatShot
	if n >= 0 && n < len(imgHeatShotAtk) {
		pos := objanim.GetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X+100, view.Y-20, 1, 0, imgHeatShotAtk[n], true)
	}
}

func (p *heatShot) Process() (bool, error) {
	if p.count == 5 {
		sound.On(sound.SEGun)

		pos := objanim.GetObjPos(p.Arg.OwnerID)
		for x := pos.X + 1; x < field.FieldNum.X; x++ {
			target := common.Point{X: x, Y: pos.Y}
			if field.GetPanelInfo(target).ObjectID != "" {
				// Hit
				damage.New(damage.Damage{
					Pos:           target,
					Power:         int(p.Arg.Power),
					TTL:           1,
					TargetType:    p.Arg.TargetType,
					HitEffectType: effect.TypeHeatHit,
				})

				// 誘爆
				target.X++
				anim.New(effect.Get(effect.TypeHeatHit, target, 0))
				damage.New(damage.Damage{
					Pos:           target,
					Power:         int(p.Arg.Power),
					TTL:           1,
					TargetType:    p.Arg.TargetType,
					HitEffectType: effect.TypeNone,
				})

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
		AnimType: anim.AnimTypeEffect,
	}
}

func (p *heatShot) StopByOwner() {
	if p.count < 5 {
		anim.Delete(p.ID)
	}
}
