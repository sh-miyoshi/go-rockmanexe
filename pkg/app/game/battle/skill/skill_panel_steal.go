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

// Note:
//   基本的にはアリアスチールと同じなので必要に応じて借りてくる
//   ただ対象位置判定のロジックが微妙に異なるのでわかりやすくするため分けた

type skillPanelSteal struct {
	ID  string
	Arg Argument

	count       int
	state       int
	target      common.Point
	myPanelType int
}

func newPanelSteal(objID string, arg Argument) *skillPanelSteal {
	res := &skillPanelSteal{
		ID:    objID,
		Arg:   arg,
		state: areaStealStateBlackout,
	}

	if arg.TargetType == field.PanelTypePlayer {
		res.myPanelType = field.PanelTypeEnemy
	} else {
		res.myPanelType = field.PanelTypePlayer
	}

	return res
}

func (p *skillPanelSteal) Draw() {
	switch p.state {
	case areaStealStateBlackout:
	case areaStealStateActing:
		ofs := p.count*4 - 30
		ino := p.count / 3
		if ino >= len(imgAreaStealMain) {
			ino = len(imgAreaStealMain) - 1
		}

		view := battlecommon.ViewPos(p.target)
		dxlib.DrawRotaGraph(view.X, view.Y+ofs, 1, 0, imgAreaStealMain[ino], true)
	case areaStealStateHit:
		ino := p.count / delayAreaStealHit
		if ino >= len(imgAreaStealPanel) {
			ino = len(imgAreaStealPanel) - 1
		}
		view := battlecommon.ViewPos(p.target)
		dxlib.DrawRotaGraph(view.X, view.Y+30, 1, 0, imgAreaStealPanel[ino], true)
	}
}

func (p *skillPanelSteal) Process() (bool, error) {
	p.count++

	switch p.state {
	case areaStealStateBlackout:
		if p.count == 1 {
			sound.On(sound.SEAreaSteal)
			field.SetBlackoutCount(90)
			setChipNameDraw("パネルスチール")

			// Target位置を実行時の一番最初に設定する
			if p.myPanelType == field.PanelTypePlayer {
				for x := 1; x < field.FieldNum.X; x++ {
					pos := objanim.GetObjPos(p.Arg.OwnerID)
					target := common.Point{X: x, Y: pos.Y}
					pn := field.GetPanelInfo(target)
					if pn.Type != field.PanelTypePlayer {
						p.target = target
						return false, nil
					}
				}
			} else {
				for x := field.FieldNum.X - 2; x >= 0; x-- {
					pos := objanim.GetObjPos(p.Arg.OwnerID)
					target := common.Point{X: x, Y: pos.Y}
					pn := field.GetPanelInfo(target)
					if pn.Type != field.PanelTypeEnemy {
						p.target = target
						return false, nil
					}
				}
			}
		}
		if p.count == 30 {
			p.setState(areaStealStateActing)
		}
	case areaStealStateActing:
		if p.count == 15 {
			sound.On(sound.SEAreaStealHit)
			p.setState(areaStealStateHit)
		}
	case areaStealStateHit:
		max := delayAreaStealHit * len(imgAreaStealPanel)
		if p.count >= max {
			pn := field.GetPanelInfo(p.target)
			if pn.ObjectID != "" {
				// ダメージ
				damage.New(damage.Damage{
					Pos:           p.target,
					Power:         10,
					TTL:           1,
					TargetType:    p.Arg.TargetType,
					HitEffectType: effect.TypeNone,
					BigDamage:     false,
					DamageType:    damage.TypeNone,
				})
			} else {
				// パネルを塗り替え
				// TODO 最終ラインの場合は塗り替えない
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
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *skillPanelSteal) StopByOwner() {
	anim.Delete(p.ID)
}

func (p *skillPanelSteal) setState(next int) {
	p.state = next
	p.count = 0
}
