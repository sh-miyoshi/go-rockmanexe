package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type spreadGun struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count int
}

type spreadHit struct {
	ID         string
	Power      uint
	TargetType int

	count int
	x, y  int
}

func (p *spreadGun) Draw() {
	n := p.count / delaySpreadGun

	// Show body
	if n < len(imgSpreadGunBody) {
		px, py := objanim.GetObjPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x+50, y-18, 1, 0, imgSpreadGunBody[n], dxlib.TRUE)
	}

	// Show atk
	n = (p.count - 4) / delaySpreadGun
	if n >= 0 && n < len(imgSpreadGunAtk) {
		px, py := objanim.GetObjPos(p.OwnerID)
		x, y := battlecommon.ViewPos(px, py)
		dxlib.DrawRotaGraph(x+100, y-20, 1, 0, imgSpreadGunAtk[n], dxlib.TRUE)
	}
}

func (p *spreadGun) Process() (bool, error) {
	if p.count == 5 {
		sound.On(sound.SEGun)

		px, py := objanim.GetObjPos(p.OwnerID)
		for x := px + 1; x < field.FieldNumX; x++ {
			if field.GetPanelInfo(x, py).ObjectID != "" {
				// Hit
				sound.On(sound.SESpreadHit)

				damage.New(damage.Damage{
					PosX:          x,
					PosY:          py,
					Power:         int(p.Power),
					TTL:           1,
					TargetType:    p.TargetType,
					HitEffectType: effect.TypeHitBig,
				})
				// Spreading
				for sy := -1; sy <= 1; sy++ {
					if py+sy < 0 || py+sy >= field.FieldNumY {
						continue
					}
					for sx := -1; sx <= 1; sx++ {
						if sy == 0 && sx == 0 {
							continue
						}
						if x+sx >= 0 && x+sx < field.FieldNumX {
							anim.New(&spreadHit{
								Power:      p.Power,
								TargetType: p.TargetType,
								x:          x + sx,
								y:          py + sy,
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
		AnimType: anim.AnimTypeEffect,
	}
}

func (p *spreadHit) Draw() {
}

func (p *spreadHit) Process() (bool, error) {
	p.count++
	if p.count == 10 {
		anim.New(effect.Get(effect.TypeSpreadHit, p.x, p.y, 5))
		damage.New(damage.Damage{
			PosX:          p.x,
			PosY:          p.y,
			Power:         int(p.Power),
			TTL:           1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
		})

		return true, nil
	}
	return false, nil
}

func (p *spreadHit) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeEffect,
	}
}
