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
	delayVulcan = 2
)

type vulcan struct {
	ID    string
	Arg   Argument
	Times int

	count    int
	imageNo  int
	atkCount int
	hit      bool
}

func newVulcan(objID string, arg Argument) *vulcan {
	return &vulcan{
		ID:    objID,
		Arg:   arg,
		Times: 3,
	}
}

func (p *vulcan) Draw() {
	pos := objanim.GetObjPos(p.Arg.OwnerID)
	view := battlecommon.ViewPos(pos)

	// Show body
	dxlib.DrawRotaGraph(view.X+50, view.Y-18, 1, 0, imgVulcan[p.imageNo], true)
	// Show attack
	if p.imageNo != 0 {
		if p.imageNo%2 == 0 {
			dxlib.DrawRotaGraph(view.X+100, view.Y-10, 1, 0, imgVulcan[3], true)
		} else {
			dxlib.DrawRotaGraph(view.X+100, view.Y-15, 1, 0, imgVulcan[3], true)
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
			pos := objanim.GetObjPos(p.Arg.OwnerID)
			hit := false
			p.atkCount++
			lastAtk := p.atkCount == p.Times
			for x := pos.X + 1; x < field.FieldNum.X; x++ {
				target := common.Point{X: x, Y: pos.Y}
				if field.GetPanelInfo(target).ObjectID != "" {
					damage.New(damage.Damage{
						Pos:           target,
						Power:         int(p.Arg.Power),
						TTL:           1,
						TargetType:    p.Arg.TargetType,
						HitEffectType: effect.TypeSpreadHit,
						BigDamage:     lastAtk,
						DamageType:    damage.TypeNone,
					})
					anim.New(effect.Get(effect.TypeVulcanHit1, target, 20))
					if p.hit && x < field.FieldNum.X-1 {
						target = common.Point{X: x + 1, Y: pos.Y}
						anim.New(effect.Get(effect.TypeVulcanHit2, target, 20))
						damage.New(damage.Damage{
							Pos:           target,
							Power:         int(p.Arg.Power),
							TTL:           1,
							TargetType:    p.Arg.TargetType,
							HitEffectType: effect.TypeNone,
							BigDamage:     lastAtk,
							DamageType:    damage.TypeNone,
						})
					}
					hit = true
					sound.On(sound.SECannonHit)
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
		AnimType: anim.AnimTypeEffect,
	}
}

func (p *vulcan) StopByOwner() {
	anim.Delete(p.ID)
}
