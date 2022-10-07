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

type heatShotHit struct {
	ID  string
	Arg Argument

	count int
	pos   common.Point
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
				sound.On(sound.SESpreadHit)

				damage.New(damage.Damage{
					Pos:           target,
					Power:         int(p.Arg.Power),
					TTL:           1,
					TargetType:    p.Arg.TargetType,
					HitEffectType: effect.TypeHitBig,
				})
				// Spreading
				for sy := -1; sy <= 1; sy++ {
					if pos.Y+sy < 0 || pos.Y+sy >= field.FieldNum.Y {
						continue
					}
					for sx := -1; sx <= 1; sx++ {
						if sy == 0 && sx == 0 {
							continue
						}
						if x+sx >= 0 && x+sx < field.FieldNum.X {
							anim.New(&heatShotHit{
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

func (p *heatShotHit) Draw() {
}

func (p *heatShotHit) Process() (bool, error) {
	p.count++
	if p.count == 10 {
		anim.New(effect.Get(effect.TypeSpreadHit, p.pos, 5))
		damage.New(damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeNone,
		})

		return true, nil
	}
	return false, nil
}

func (p *heatShotHit) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeEffect,
	}
}

func (p *heatShotHit) StopByOwner() {
	anim.Delete(p.ID)
}
