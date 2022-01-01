package skill

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	thunderBallNextStepCount = 80
)

type thunderBall struct {
	ID           string
	OwnerID      string
	Power        uint
	TargetType   int
	MaxMoveCount int

	count     int
	x, y      int
	moveCount int
	damageID  string
	nextX     int
	nextY     int
	prevX     int
	prevY     int
}

func (p *thunderBall) Draw() {
	x, y := battlecommon.ViewPos(p.x, p.y)
	n := (p.count / delayThunderBall) % len(imgThunderBall)

	cnt := p.count % thunderBallNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(p.nextX, p.x, p.prevX, cnt, thunderBallNextStepCount, field.PanelSizeX)
	ofsy := battlecommon.GetOffset(p.nextY, p.y, p.prevY, cnt, thunderBallNextStepCount, field.PanelSizeY)

	dxlib.DrawRotaGraph(x+int32(ofsx), y+25+int32(ofsy), 1, 0, imgThunderBall[n], dxlib.TRUE)
}

func (p *thunderBall) Process() (bool, error) {
	if p.count == 0 {
		sound.On(sound.SEThunderBall)
	}

	if p.count%thunderBallNextStepCount == 2 {
		if p.damageID != "" {
			if !damage.Exists(p.damageID) {
				// attack hit to target
				return true, nil
			}
		}
	}

	if p.count%thunderBallNextStepCount == 0 {
		tx := p.x
		ty := p.y
		if p.count != 0 {
			// Update current pos
			p.prevX = p.x
			p.prevY = p.y
			p.x = p.nextX
			p.y = p.nextY

			p.moveCount++
			if p.moveCount > p.MaxMoveCount {
				return true, nil
			}

			if p.x < 0 || p.x > field.FieldNumX || p.y < 0 || p.y > field.FieldNumY {
				return true, nil
			}
		}

		pn := field.GetPanelInfo(p.x, p.y)
		if pn.Status == field.PanelStatusHole {
			return true, nil
		}

		p.damageID = damage.New(damage.Damage{
			PosX:          p.x,
			PosY:          p.y,
			Power:         int(p.Power),
			TTL:           thunderBallNextStepCount + 1,
			TargetType:    p.TargetType,
			HitEffectType: effect.TypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
		})

		// Set next pos
		objType := objanim.ObjTypePlayer
		if p.TargetType == damage.TargetEnemy {
			objType = objanim.ObjTypeEnemy
		}

		objs := objanim.GetObjs(objanim.Filter{ObjType: objType})
		if len(objs) == 0 {
			// no target
			if p.TargetType == damage.TargetPlayer {
				p.nextX--
			} else {
				p.nextX++
			}
		} else {
			xdif := objs[0].PosX - tx
			ydif := objs[0].PosY - ty

			if xdif != 0 || ydif != 0 {
				if common.Abs(xdif) > common.Abs(ydif) {
					// move to x
					p.nextX += (xdif / common.Abs(xdif))
				} else {
					// move to y
					p.nextY += (ydif / common.Abs(ydif))
				}
			}
		}
	}

	p.count++
	return false, nil
}

func (p *thunderBall) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}
