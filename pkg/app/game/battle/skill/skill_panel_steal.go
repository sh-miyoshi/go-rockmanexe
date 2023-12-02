package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

// Note:
//   基本的にはエリアスチールと同じなので必要に応じて借りてくる
//   ただ対象位置判定のロジックが微妙に異なるのでわかりやすくするため分けた

type skillPanelSteal struct {
	ID  string
	Arg Argument

	count       int
	state       int
	target      point.Point
	myPanelType int
	drawer      skilldraw.DrawAreaSteal
}

func newPanelSteal(objID string, arg Argument) *skillPanelSteal {
	res := &skillPanelSteal{
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

func (p *skillPanelSteal) Draw() {
	p.drawer.Draw(p.count, p.state, []point.Point{p.target})
}

func (p *skillPanelSteal) Process() (bool, error) {
	p.count++

	switch p.state {
	case resources.SkillAreaStealStateBlackout:
		if p.count == 1 {
			sound.On(resources.SEAreaSteal)
			field.SetBlackoutCount(90)
			SetChipNameDraw("パネルスチール", true)

			// Target位置を実行時の一番最初に設定する
			if p.myPanelType == battlecommon.PanelTypePlayer {
				for x := 1; x < battlecommon.FieldNum.X; x++ {
					pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
					target := point.Point{X: x, Y: pos.Y}
					pn := field.GetPanelInfo(target)
					if pn.Type != battlecommon.PanelTypePlayer {
						p.target = target
						return false, nil
					}
				}
			} else {
				for x := battlecommon.FieldNum.X - 2; x >= 0; x-- {
					pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
					target := point.Point{X: x, Y: pos.Y}
					pn := field.GetPanelInfo(target)
					if pn.Type != battlecommon.PanelTypeEnemy {
						p.target = target
						return false, nil
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
			pn := field.GetPanelInfo(p.target)
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
			} else if p.target.X >= 1 && p.target.X < battlecommon.FieldNum.X-1 {
				// パネルを塗り替え
				// 最終ラインの場合は塗り替えない
				field.ChangePanelType(p.target, p.myPanelType)
			}
			return true, nil
		}
	}

	return false, nil
}

func (p *skillPanelSteal) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillPanelSteal) StopByOwner() {
	localanim.AnimDelete(p.ID)
}

func (p *skillPanelSteal) setState(next int) {
	p.state = next
	p.count = 0
}
