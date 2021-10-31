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

type vulcan struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int
	Times      int

	count    int
	imageNo  int
	atkCount int
	hit      bool
}

func (p *vulcan) Draw() {
	px, py := objanim.GetObjPos(p.OwnerID)
	x, y := battlecommon.ViewPos(px, py)

	// Show body
	dxlib.DrawRotaGraph(x+50, y-18, 1, 0, imgVulcan[p.imageNo], dxlib.TRUE)
	// Show attack
	if p.imageNo != 0 {
		if p.imageNo%2 == 0 {
			dxlib.DrawRotaGraph(x+100, y-10, 1, 0, imgVulcan[3], dxlib.TRUE)
		} else {
			dxlib.DrawRotaGraph(x+100, y-15, 1, 0, imgVulcan[3], dxlib.TRUE)
		}
	}
}

func (p *vulcan) Process() (bool, error) {
	p.count++
	if p.count >= delayVulcan*1 {
		if p.count%(delayVulcan*5) == delayVulcan*1 {
			sound.On(sound.SEGun)

			p.imageNo = p.imageNo%2 + 1
			// Add damage
			px, py := objanim.GetObjPos(p.OwnerID)
			hit := false
			for x := px + 1; x < field.FieldNumX; x++ {
				if field.GetPanelInfo(x, py).ObjectID != "" {
					damage.New(damage.Damage{
						PosX:          x,
						PosY:          py,
						Power:         int(p.Power),
						TTL:           1,
						TargetType:    p.TargetType,
						HitEffectType: effect.TypeSpreadHit,
					})
					anim.New(effect.Get(effect.TypeVulcanHit1, x, py, 20))
					if p.hit && x < field.FieldNumX-1 {
						anim.New(effect.Get(effect.TypeVulcanHit2, x+1, py, 20))
						damage.New(damage.Damage{
							PosX:          x + 1,
							PosY:          py,
							Power:         int(p.Power),
							TTL:           1,
							TargetType:    p.TargetType,
							HitEffectType: effect.TypeNone,
						})
					}
					hit = true
					sound.On(sound.SECannonHit)
					break
				}
			}
			p.hit = hit
			p.atkCount++
			if p.atkCount == p.Times {
				return true, nil
			}
		}

	}

	return false, nil
}

func (p *vulcan) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeEffect,
	}
}
