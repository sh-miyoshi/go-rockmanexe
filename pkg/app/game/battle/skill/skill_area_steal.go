package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

type skillAreaSteal struct {
	ID  string
	Arg Argument

	count       int
	state       int
	targetLineX int
	myPanelType int
	drawer      skilldraw.DrawAreaSteal
}

func newAreaSteal(objID string, arg Argument) *skillAreaSteal {
	res := &skillAreaSteal{
		ID:    objID,
		Arg:   arg,
		state: resources.SkillAreaStealStateBlackout,
	}

	if arg.TargetType == battlecommon.PanelTypePlayer {
		res.myPanelType = battlecommon.PanelTypeEnemy
	} else {
		res.myPanelType = battlecommon.PanelTypePlayer
	}

	return res
}

func (p *skillAreaSteal) Draw() {
	targets := []common.Point{}
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		targets = append(targets, common.Point{X: p.targetLineX, Y: y})
	}

	p.drawer.Draw(p.count, p.state, targets)
}

func (p *skillAreaSteal) Process() (bool, error) {
	p.count++

	switch p.state {
	case resources.SkillAreaStealStateBlackout:
		if p.count == 1 {
			sound.On(resources.SEAreaSteal)
			field.SetBlackoutCount(90)
			SetChipNameDraw("エリアスチール", true)

			// Target Lineを実行時の一番最初に設定する
			if p.myPanelType == battlecommon.PanelTypePlayer {
				for x := 1; x < battlecommon.FieldNum.X; x++ {
					for y := 0; y < battlecommon.FieldNum.Y; y++ {
						pn := field.GetPanelInfo(common.Point{X: x, Y: y})
						if pn.Type != battlecommon.PanelTypePlayer {
							p.targetLineX = x
							return false, nil
						}
					}
				}
			} else {
				for x := battlecommon.FieldNum.X - 2; x >= 0; x-- {
					for y := 0; y < battlecommon.FieldNum.Y; y++ {
						pn := field.GetPanelInfo(common.Point{X: x, Y: y})
						if pn.Type != battlecommon.PanelTypeEnemy {
							p.targetLineX = x
							return false, nil
						}
					}
				}
			}
		}
		if p.count == 30 {
			p.setState(resources.SkillAreaStealStateActing)
		}
	case resources.SkillAreaStealStateActing:
		if p.count == 15 {
			sound.On(resources.SEAreaStealHit)
			p.setState(resources.SkillAreaStealStateHit)
		}
	case resources.SkillAreaStealStateHit:
		if p.count >= resources.SkillAreaStealHitEndCount {
			for y := 0; y < battlecommon.FieldNum.Y; y++ {
				pos := common.Point{X: p.targetLineX, Y: y}
				pn := field.GetPanelInfo(pos)
				if pn.ObjectID != "" {
					// ダメージ
					localanim.DamageManager().New(damage.Damage{
						DamageType:    damage.TypeObject,
						Power:         10,
						TargetObjType: p.Arg.TargetType,
						HitEffectType: resources.EffectTypeNone,
						BigDamage:     false,
						Element:       damage.ElementNone,
						TargetObjID:   pn.ObjectID,
					})
				} else if p.targetLineX >= 1 && p.targetLineX < battlecommon.FieldNum.X-1 {
					// パネルを塗り替え
					// 最終ラインの場合は塗り替えない
					field.ChangePanelType(pos, p.myPanelType)
				}
			}
			return true, nil
		}
	}

	return false, nil
}

func (p *skillAreaSteal) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillAreaSteal) StopByOwner() {
	localanim.AnimDelete(p.ID)
}

func (p *skillAreaSteal) setState(next int) {
	p.state = next
	p.count = 0
}
