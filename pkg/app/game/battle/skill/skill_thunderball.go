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
	thunderBallNextStepCount = 80
	delayThunderBall         = 6
)

type thunderBall struct {
	ID           string
	Arg          Argument
	MaxMoveCount int

	count     int
	moveCount int
	damageID  string
	pos       common.Point
	next      common.Point
	prev      common.Point
}

func newThunderBall(objID string, arg Argument) *thunderBall {
	pos := objanim.GetObjPos(arg.OwnerID)
	x := pos.X + 1
	if arg.TargetType == damage.TargetPlayer {
		x = pos.X - 1
	}

	first := common.Point{X: x, Y: pos.Y}
	max := 6 // debug
	return &thunderBall{
		ID:           objID,
		Arg:          arg,
		MaxMoveCount: max,
		pos:          first,
		prev:         pos,
		next:         first,
	}
}

func (p *thunderBall) Draw() {
	view := battlecommon.ViewPos(p.pos)
	n := (p.count / delayThunderBall) % len(imgThunderBall)

	cnt := p.count % thunderBallNextStepCount
	if cnt == 0 {
		// Skip drawing because the position is updated in Process method and return unexpected value
		return
	}

	ofsx := battlecommon.GetOffset(p.next.X, p.pos.X, p.prev.X, cnt, thunderBallNextStepCount, field.PanelSize.X)
	ofsy := battlecommon.GetOffset(p.next.Y, p.pos.Y, p.prev.Y, cnt, thunderBallNextStepCount, field.PanelSize.Y)

	dxlib.DrawRotaGraph(view.X+ofsx, view.Y+25+ofsy, 1, 0, imgThunderBall[n], true)
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
		t := p.pos
		if p.count != 0 {
			// Update current pos
			p.prev = p.pos
			p.pos = p.next

			p.moveCount++
			if p.moveCount > p.MaxMoveCount {
				return true, nil
			}

			if p.pos.X < 0 || p.pos.X > field.FieldNum.X || p.pos.Y < 0 || p.pos.Y > field.FieldNum.Y {
				return true, nil
			}
		}

		pn := field.GetPanelInfo(p.pos)
		if pn.Status == field.PanelStatusHole {
			return true, nil
		}

		p.damageID = damage.New(damage.Damage{
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           thunderBallNextStepCount + 1,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
			DamageType:    damage.TypeElec,
		})

		// Set next pos
		objType := objanim.ObjTypePlayer
		if p.Arg.TargetType == damage.TargetEnemy {
			objType = objanim.ObjTypeEnemy
		}

		objs := objanim.GetObjs(objanim.Filter{ObjType: objType})
		if len(objs) == 0 {
			// no target
			if p.Arg.TargetType == damage.TargetPlayer {
				p.next.X--
			} else {
				p.next.X++
			}
		} else {
			xdif := objs[0].Pos.X - t.X
			ydif := objs[0].Pos.Y - t.Y

			if xdif != 0 || ydif != 0 {
				if common.Abs(xdif) > common.Abs(ydif) {
					// move to x
					p.next.X += (xdif / common.Abs(xdif))
				} else {
					// move to y
					p.next.Y += (ydif / common.Abs(ydif))
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

func (p *thunderBall) StopByOwner() {
	// Nothing to do
}
