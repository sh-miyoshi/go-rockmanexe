package processor

import (
	"github.com/google/uuid"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	AirHockeyNextStepCount = 3
	airHockeyMoveCountMax  = 9
)

type AirHockey struct {
	SkillID int
	Arg     skillcore.Argument

	count     int
	moveCount int
	pos       point.Point
	next      point.Point
	prev      point.Point
	moveVec   point.Point
}

func (p *AirHockey) Init() {
	pos := p.Arg.GetObjectPos(p.Arg.OwnerID)
	x := pos.X + 1
	if p.Arg.TargetType == damage.TargetPlayer {
		x = pos.X - 1
	}

	first := point.Point{X: x, Y: pos.Y}
	p.pos = first
	p.prev = pos
	p.next = first

	// プレイヤーの場合は右下、敵の場合は左下に移動
	if p.Arg.TargetType == damage.TargetPlayer {
		p.moveVec = point.Point{X: -1, Y: 1}
	} else {
		p.moveVec = point.Point{X: 1, Y: 1}
	}
}

func (p *AirHockey) Update() (bool, error) {
	if p.moveCount >= airHockeyMoveCountMax {
		return true, nil
	}

	// AirHockeyNextStepCountごとに移動
	if p.count%AirHockeyNextStepCount == 0 {
		nextX := p.pos.X + p.moveVec.X
		nextY := p.pos.Y + p.moveVec.Y

		// 画面外への移動時の反射処理
		if nextX < 0 || nextX > battlecommon.FieldNum.X-1 {
			p.moveVec.X = -p.moveVec.X
			nextX = p.pos.X + p.moveVec.X
		}
		if nextY < 0 || nextY > battlecommon.FieldNum.Y-1 {
			p.moveVec.Y = -p.moveVec.Y
			nextY = p.pos.Y + p.moveVec.Y
		}

		// エリアの端での反射処理
		currentPanel := p.Arg.GetPanelInfo(p.pos)
		nextPanel := p.Arg.GetPanelInfo(point.Point{X: nextX, Y: nextY})
		if p.Arg.TargetType == damage.TargetEnemy && currentPanel.Type == battlecommon.PanelTypeEnemy && nextPanel.Type == battlecommon.PanelTypePlayer {
			// 敵エリアからプレイヤーエリアへの移動時は反射
			p.moveVec.X = -p.moveVec.X
			nextX = p.pos.X + p.moveVec.X
		} else if p.Arg.TargetType == damage.TargetPlayer && currentPanel.Type == battlecommon.PanelTypePlayer && nextPanel.Type == battlecommon.PanelTypeEnemy {
			// プレイヤーエリアから敵エリアへの移動時は反射
			p.moveVec.X = -p.moveVec.X
			nextX = p.pos.X + p.moveVec.X
		}

		// 移動前の位置を保存
		p.prev = p.pos
		// 次の位置を計算
		next := point.Point{X: nextX, Y: nextY}

		// 同じ位置への移動の場合は終了
		if p.pos.X == next.X && p.pos.Y == next.Y {
			return true, nil
		}

		// ダメージ処理
		if objID := p.Arg.GetPanelInfo(p.pos).ObjectID; objID != "" {
			dm := damage.Damage{
				OwnerClientID: p.Arg.OwnerClientID,
				ID:            uuid.New().String(),
				Power:         int(p.Arg.Power),
				HitEffectType: resources.EffectTypeNone,
				StrengthType:  damage.StrengthBack,
				Element:       damage.ElementNone, // WIP: ブレイク属性にする

				DamageType:    damage.TypeObject,
				TargetObjType: p.Arg.TargetType,
				TargetObjID:   objID,
			}
			p.Arg.DamageMgr.New(dm)
		}

		p.pos = p.next
		p.next = next
		p.moveCount++
	}

	p.count++
	return false, nil
}

func (p *AirHockey) GetCount() int {
	return p.count
}

func (p *AirHockey) GetPos() (prev, current, next point.Point) {
	return p.prev, p.pos, p.next
}
