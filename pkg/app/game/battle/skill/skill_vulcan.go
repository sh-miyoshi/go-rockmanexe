package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
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
	pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
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
			pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
			hit := false
			p.atkCount++
			lastAtk := p.atkCount == p.Times
			for x := pos.X + 1; x < battlecommon.FieldNum.X; x++ {
				target := common.Point{X: x, Y: pos.Y}
				if field.GetPanelInfo(target).ObjectID != "" {
					localanim.DamageManager().New(damage.Damage{
						Pos:           target,
						Power:         int(p.Arg.Power),
						TTL:           1,
						TargetType:    p.Arg.TargetType,
						HitEffectType: battlecommon.EffectTypeSpreadHit,
						BigDamage:     lastAtk,
						DamageType:    damage.TypeNone,
					})
					localanim.AnimNew(effect.Get(battlecommon.EffectTypeVulcanHit1, target, 20))
					if p.hit && x < battlecommon.FieldNum.X-1 {
						target = common.Point{X: x + 1, Y: pos.Y}
						localanim.AnimNew(effect.Get(battlecommon.EffectTypeVulcanHit2, target, 20))
						localanim.DamageManager().New(damage.Damage{
							Pos:           target,
							Power:         int(p.Arg.Power),
							TTL:           1,
							TargetType:    p.Arg.TargetType,
							HitEffectType: battlecommon.EffectTypeNone,
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
		DrawType: anim.DrawTypeEffect,
	}
}

func (p *vulcan) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
