package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayAreaStealHit = 2
)

const (
	areaStealStateBlackout int = iota
	areaStealStateActing
	areaStealStateHit
)

type skillAreaSteal struct {
	ID  string
	Arg Argument

	count       int
	state       int
	targetLineX int
	myPanelType int
}

func newAreaSteal(objID string, arg Argument) *skillAreaSteal {
	res := &skillAreaSteal{
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

func (p *skillAreaSteal) Draw() {
	switch p.state {
	case areaStealStateBlackout:
	case areaStealStateActing:
		ofs := p.count*4 - 30
		ino := p.count / 3
		if ino >= len(imgAreaStealMain) {
			ino = len(imgAreaStealMain) - 1
		}

		for y := 0; y < field.FieldNum.Y; y++ {
			view := battlecommon.ViewPos(common.Point{X: p.targetLineX, Y: y})
			dxlib.DrawRotaGraph(view.X, view.Y+ofs, 1, 0, imgAreaStealMain[ino], true)
		}
	case areaStealStateHit:
		ino := p.count / delayAreaStealHit
		if ino >= len(imgAreaStealPanel) {
			ino = len(imgAreaStealPanel) - 1
		}
		for y := 0; y < field.FieldNum.Y; y++ {
			view := battlecommon.ViewPos(common.Point{X: p.targetLineX, Y: y})
			dxlib.DrawRotaGraph(view.X, view.Y+30, 1, 0, imgAreaStealPanel[ino], true)
		}
	}
}

func (p *skillAreaSteal) Process() (bool, error) {
	p.count++

	switch p.state {
	case areaStealStateBlackout:
		if p.count == 1 {
			sound.On(sound.SEAreaSteal)
			field.SetBlackoutCount(90)
			setChipNameDraw("エリアスチール")

			// Target Lineを実行時の一番最初に設定する
			if p.myPanelType == field.PanelTypePlayer {
				for x := 1; x < field.FieldNum.X; x++ {
					for y := 0; y < field.FieldNum.Y; y++ {
						pn := field.GetPanelInfo(common.Point{X: x, Y: y})
						if pn.Type != field.PanelTypePlayer {
							p.targetLineX = x
							return false, nil
						}
					}
				}
			} else {
				for x := field.FieldNum.X - 2; x >= 0; x-- {
					for y := 0; y < field.FieldNum.Y; y++ {
						pn := field.GetPanelInfo(common.Point{X: x, Y: y})
						if pn.Type != field.PanelTypeEnemy {
							p.targetLineX = x
							return false, nil
						}
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
			for y := 0; y < field.FieldNum.Y; y++ {
				pos := common.Point{X: p.targetLineX, Y: y}
				pn := field.GetPanelInfo(pos)
				if pn.ObjectID != "" {
					// ダメージ
					damage.New(damage.Damage{
						Pos:           pos,
						Power:         10,
						TTL:           1,
						TargetType:    p.Arg.TargetType,
						HitEffectType: effect.TypeNone,
						BigDamage:     false,
						DamageType:    damage.TypeNone,
					})
				} else if p.targetLineX >= 1 && p.targetLineX < field.FieldNum.X-1 {
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
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *skillAreaSteal) StopByOwner() {
	anim.Delete(p.ID)
}

func (p *skillAreaSteal) setState(next int) {
	p.state = next
	p.count = 0
}
