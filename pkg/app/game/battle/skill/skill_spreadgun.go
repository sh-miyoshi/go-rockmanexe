package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delaySpreadGun = 2
)

type spreadGun struct {
	ID  string
	Arg Argument

	count int
}

type spreadHit struct {
	ID  string
	Arg Argument

	count int
	pos   common.Point
}

func newSpreadGun(objID string, arg Argument) *spreadGun {
	return &spreadGun{
		ID:  objID,
		Arg: arg,
	}
}

func (p *spreadGun) Draw() {
	n := p.count / delaySpreadGun

	// Show body
	if n < len(imgSpreadGunBody) {
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X+50, view.Y-18, 1, 0, imgSpreadGunBody[n], true)
	}

	// Show atk
	n = (p.count - 4) / delaySpreadGun
	if n >= 0 && n < len(imgSpreadGunAtk) {
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		view := battlecommon.ViewPos(pos)
		dxlib.DrawRotaGraph(view.X+100, view.Y-20, 1, 0, imgSpreadGunAtk[n], true)
	}
}

func (p *spreadGun) Process() (bool, error) {
	if p.count == 5 {
		sound.On(sound.SEGun)

		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
			target := common.Point{X: x, Y: pos.Y}
			objs := localanim.ObjAnimGetObjs(objanim.Filter{Pos: &target, ObjType: p.Arg.TargetType})
			if len(objs) > 0 {
				// Hit
				sound.On(sound.SESpreadHit)

				localanim.DamageManager().New(damage.Damage{
					Pos:           target,
					Power:         int(p.Arg.Power),
					TTL:           1,
					TargetType:    p.Arg.TargetType,
					HitEffectType: battlecommon.EffectTypeHitBig,
					DamageType:    damage.TypeNone,
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

	max := len(imgSpreadGunAtk)
	if len(imgSpreadGunBody) > max {
		max = len(imgSpreadGunBody)
	}

	if p.count > max*delaySpreadGun {
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
		localanim.AnimNew(effect.Get(battlecommon.EffectTypeSpreadHit, p.pos, 5))
		localanim.DamageManager().New(damage.Damage{
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
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *spreadHit) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
