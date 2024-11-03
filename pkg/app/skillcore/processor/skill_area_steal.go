package processor

import (
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	areaStealHitEndCount = 12
)

type AreaSteal struct {
	Arg skillcore.Argument

	count       int
	state       int
	targets     []point.Point
	myPanelType int
	name        string
}

func (p *AreaSteal) Init(skillID int) {
	p.setState(resources.SkillAreaStealStateBlackout)

	if p.Arg.TargetType == battlecommon.PanelTypePlayer {
		p.myPanelType = battlecommon.PanelTypeEnemy
	} else {
		p.myPanelType = battlecommon.PanelTypePlayer
	}

	if skillID == resources.SkillAreaSteal {
		p.name = "エリアスチール"
		lineX := -1
		if p.myPanelType == battlecommon.PanelTypePlayer {
		FindLoopPlayer:
			for x := 1; x < battlecommon.FieldNum.X; x++ {
				for y := 0; y < battlecommon.FieldNum.Y; y++ {
					pn := p.Arg.GetPanelInfo(point.Point{X: x, Y: y})
					if pn.Type != battlecommon.PanelTypePlayer {
						lineX = x
						break FindLoopPlayer
					}
				}
			}
		} else {
		FindLoopEnemy:
			for x := battlecommon.FieldNum.X - 2; x >= 0; x-- {
				for y := 0; y < battlecommon.FieldNum.Y; y++ {
					pn := p.Arg.GetPanelInfo(point.Point{X: x, Y: y})
					if pn.Type != battlecommon.PanelTypeEnemy {
						lineX = x
						break FindLoopEnemy
					}
				}
			}
		}
		if lineX != -1 {
			p.targets = []point.Point{}
			for y := 0; y < battlecommon.FieldNum.Y; y++ {
				p.targets = append(p.targets, point.Point{X: lineX, Y: y})
			}
		}
	} else {
		p.name = "パネルスチール"
		if p.myPanelType == battlecommon.PanelTypePlayer {
			for x := 1; x < battlecommon.FieldNum.X; x++ {
				pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
				target := point.Point{X: x, Y: pos.Y}
				pn := p.Arg.GetPanelInfo(target)
				if pn.Type != battlecommon.PanelTypePlayer {
					p.targets = []point.Point{target}
					return
				}
			}
		} else {
			for x := battlecommon.FieldNum.X - 2; x >= 0; x-- {
				pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
				target := point.Point{X: x, Y: pos.Y}
				pn := p.Arg.GetPanelInfo(target)
				if pn.Type != battlecommon.PanelTypeEnemy {
					p.targets = []point.Point{target}
					return
				}
			}
		}
	}
}

func (p *AreaSteal) Update() (bool, error) {
	p.count++

	switch p.state {
	case resources.SkillAreaStealStateBlackout:
		if p.count == 1 {
			p.Arg.SoundOn(resources.SEAreaSteal)
			p.Arg.Cutin(p.name, 90)
		}

		if p.count == 30 {
			p.setState(resources.SkillAreaStealStateActing)
		}
	case resources.SkillAreaStealStateActing:
		if p.count == 15 {
			p.Arg.SoundOn(resources.SEAreaStealHit)
			p.setState(resources.SkillAreaStealStateHit)
		}
	case resources.SkillAreaStealStateHit:
		if p.count >= areaStealHitEndCount {
			for _, pos := range p.targets {
				pn := p.Arg.GetPanelInfo(pos)
				if pn.ObjectID != "" {
					// ダメージ
					p.Arg.DamageMgr.New(damage.Damage{
						DamageType:    damage.TypeObject,
						Power:         10,
						TargetObjType: p.Arg.TargetType,
						HitEffectType: resources.EffectTypeNone,
						BigDamage:     false,
						Element:       damage.ElementNone,
						TargetObjID:   pn.ObjectID,
					})
				} else if pos.X >= 1 && pos.X < battlecommon.FieldNum.X-1 {
					// パネルを塗り替え
					// 最終ラインの場合は塗り替えない
					p.Arg.ChangePanelType(pos, p.myPanelType, battlecommon.DefaultPanelTypeEndCount)
				}
			}
			return true, nil
		}
	}

	return false, nil
}

func (p *AreaSteal) GetCount() int {
	return p.count
}

func (p *AreaSteal) GetState() int {
	return p.state
}

func (p *AreaSteal) GetTargets() []point.Point {
	return p.targets
}

func (p *AreaSteal) setState(next int) {
	p.state = next
	p.count = 0
}
