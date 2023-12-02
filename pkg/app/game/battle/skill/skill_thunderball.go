package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type thunderBall struct {
	ID           string
	Arg          Argument
	MaxMoveCount int

	count     int
	moveCount int
	damageID  string
	pos       point.Point
	next      point.Point
	prev      point.Point
	drawer    skilldraw.DrawThunderBall
}

func newThunderBall(objID string, arg Argument) *thunderBall {
	pos := localanim.ObjAnimGetObjPos(arg.OwnerID)
	x := pos.X + 1
	if arg.TargetType == damage.TargetPlayer {
		x = pos.X - 1
	}

	first := point.Point{X: x, Y: pos.Y}
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
	p.drawer.Draw(p.prev, p.pos, p.next, p.count)
}

func (p *thunderBall) Process() (bool, error) {
	if p.count == 0 {
		sound.On(resources.SEThunderBall)
	}

	if p.count%resources.SkillThunderBallNextStepCount == 2 {
		if p.damageID != "" {
			if !localanim.DamageManager().Exists(p.damageID) {
				// attack hit to target
				return true, nil
			}
		}
	}

	if p.count%resources.SkillThunderBallNextStepCount == 0 {
		t := p.pos
		if p.count != 0 {
			// Update current pos
			p.prev = p.pos
			p.pos = p.next

			p.moveCount++
			if p.moveCount > p.MaxMoveCount {
				return true, nil
			}

			if p.pos.X < 0 || p.pos.X > battlecommon.FieldNum.X || p.pos.Y < 0 || p.pos.Y > battlecommon.FieldNum.Y {
				return true, nil
			}
		}

		pn := field.GetPanelInfo(p.pos)
		if pn.Status == battlecommon.PanelStatusHole {
			return true, nil
		}

		p.damageID = localanim.DamageManager().New(damage.Damage{
			DamageType:    damage.TypePosition,
			Pos:           p.pos,
			Power:         int(p.Arg.Power),
			TTL:           resources.SkillThunderBallNextStepCount + 1,
			TargetObjType: p.Arg.TargetType,
			HitEffectType: resources.EffectTypeNone,
			ShowHitArea:   true,
			BigDamage:     true,
			Element:       damage.ElementElec,
			IsParalyzed:   true,
		})

		// Set next pos
		objType := objanim.ObjTypePlayer
		if p.Arg.TargetType == damage.TargetEnemy {
			objType = objanim.ObjTypeEnemy
		}

		objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objType})
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
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *thunderBall) StopByOwner() {
	// Nothing to do
}
