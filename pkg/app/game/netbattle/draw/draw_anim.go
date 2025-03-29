package draw

import (
	"fmt"

	"github.com/cockroachdb/errors"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type animDraw struct {
	drawCannonInst   skilldraw.DrawCannon
	drawMiniBombInst skilldraw.DrawMiniBomb
	drawRecover      skilldraw.DrawRecover
	drawShockWave    skilldraw.DrawShockWave
	drawSpreadGun    skilldraw.DrawSpreadGun
	drawSpreadHit    skilldraw.DrawSpreadHit
	drawSword        skilldraw.DrawSword
	drawVulcan       skilldraw.DrawVulcan
	drawWideShot     skilldraw.DrawWideShot
	drawHeatShot     skilldraw.DrawHeatShot
	drawFlameLine    skilldraw.DrawFlamePillerManager
	drawTornado      skilldraw.DrawTornado
	drawBoomerang    skilldraw.DrawBoomerang
	drawBamboolance  skilldraw.DrawBamboolance
	drawAreaSteal    skilldraw.DrawAreaSteal
	drawBubbleShot   skilldraw.DrawBubbleShot
	drawAirHockey    skilldraw.DrawAirHockey
}

func (d *animDraw) Init() error {
	if err := skilldraw.LoadImages(); err != nil {
		return errors.Wrap(err, "failed to load skill image")
	}
	d.drawBamboolance.Init()

	return nil
}

func (d *animDraw) End() {
	skilldraw.ClearImages()
}

func (d *animDraw) Draw() {
	ginfo := net.GetInst().GetGameInfo()
	for _, a := range ginfo.Anims {
		pos := battlecommon.ViewPos(a.Pos)
		isPlayer := ginfo.ClientID == a.OwnerClientID

		switch a.AnimType {
		case anim.TypeCannonNormal, anim.TypeCannonHigh, anim.TypeCannonMega:
			skillID := resources.SkillCannon
			switch a.AnimType {
			case anim.TypeCannonHigh:
				skillID = resources.SkillHighCannon
			case anim.TypeCannonMega:
				skillID = resources.SkillMegaCannon
			}

			d.drawCannonInst.Draw(skillID, pos, a.ActCount, isPlayer)
		case anim.TypeMiniBomb:
			var drawPm skill.MiniBombDrawParam
			drawPm.Unmarshal(a.DrawParam)
			if !isPlayer {
				drawPm.Target.X = battlecommon.FieldNum.X - drawPm.Target.X - 1
			}

			d.drawMiniBombInst.Draw(a.Pos, drawPm.Target, a.ActCount, drawPm.LandCount)
		case anim.TypeRecover:
			d.drawRecover.Draw(pos, a.ActCount)
		case anim.TypeShockWave:
			var drawPm skill.ShockWaveDrawParam
			drawPm.Unmarshal(a.DrawParam)
			if a.ActCount >= drawPm.InitWait {
				if !isPlayer {
					drawPm.Direct = battlecommon.ReverseDirect(drawPm.Direct)
				}
				d.drawShockWave.Draw(pos, a.ActCount, drawPm.Speed, drawPm.Direct)
			}
		case anim.TypeSpreadGun:
			d.drawSpreadGun.Draw(pos, a.ActCount, isPlayer)
		case anim.TypeSpreadHit:
			d.drawSpreadHit.Draw(pos, a.ActCount)
		case anim.TypeSword, anim.TypeWideSword, anim.TypeLongSword:
			var drawPm skill.SwordDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawSword.Draw(drawPm.SkillID, pos, a.ActCount, drawPm.Delay, !isPlayer)
		case anim.TypeVulcan:
			var drawPm skill.VulcanDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawVulcan.Draw(pos, a.ActCount, drawPm.Delay, isPlayer)
		case anim.TypeWideShot:
			var drawPm skill.WideShotDrawParam
			drawPm.Unmarshal(a.DrawParam)
			if !isPlayer {
				drawPm.Direct = battlecommon.ReverseDirect(drawPm.Direct)
			}
			d.drawWideShot.Draw(a.Pos, a.ActCount, drawPm.Direct, true, drawPm.NextStepCount, drawPm.State)
		case anim.TypeHeatShot:
			d.drawHeatShot.Draw(pos, a.ActCount, isPlayer)
		case anim.TypeFlameLine:
			var drawPm skill.FlameLineDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawFlameLine.Draw(pos, a.ActCount, true, drawPm.Pillars, drawPm.Delay, isPlayer)
		case anim.TypeTornado:
			// Note: DrawParamで渡すようにしてもいいが、targetの決定アルゴリズムが変わることはないのでここに直接書く
			targetPos := point.Point{X: a.Pos.X + 2, Y: a.Pos.Y}
			if !isPlayer {
				targetPos.X = a.Pos.X - 2
			}
			target := battlecommon.ViewPos(targetPos)
			d.drawTornado.Draw(pos, target, a.ActCount, isPlayer)
		case anim.TypeBoomerang:
			var drawPm skill.BoomerangDrawParam
			drawPm.Unmarshal(a.DrawParam)
			if !isPlayer {
				drawPm.PrevPos.X = battlecommon.FieldNum.X - drawPm.PrevPos.X - 1
				drawPm.NextPos.X = battlecommon.FieldNum.X - drawPm.NextPos.X - 1
			}
			d.drawBoomerang.Draw(drawPm.PrevPos, a.Pos, drawPm.NextPos, a.ActCount, processor.BoomerangNextStepCount)
		case anim.TypeBambooLance:
			d.drawBamboolance.Draw(a.ActCount, isPlayer)
		case anim.TypeAreaSteal:
			var drawPm skill.AreaStealDrawParam
			drawPm.Unmarshal(a.DrawParam)
			if !isPlayer {
				for i := 0; i < len(drawPm.Targets); i++ {
					drawPm.Targets[i].X = battlecommon.FieldNum.X - drawPm.Targets[i].X - 1
				}
			}

			d.drawAreaSteal.Draw(a.ActCount, drawPm.State, drawPm.Targets)
		case anim.TypeCrack:
			// no animation
		case anim.TypeBubbleShot:
			d.drawBubbleShot.Draw(pos, a.ActCount, isPlayer)
		case anim.TypeInvisible:
			// no animation
		case anim.TypeAirHockey:
			var drawPm skill.AirHockeyDrawParam
			drawPm.Unmarshal(a.DrawParam)
			if !isPlayer {
				drawPm.PrevPos.X = battlecommon.FieldNum.X - drawPm.PrevPos.X - 1
				drawPm.NextPos.X = battlecommon.FieldNum.X - drawPm.NextPos.X - 1
			}
			d.drawAirHockey.Draw(drawPm.PrevPos, a.Pos, drawPm.NextPos, a.ActCount, processor.AirHockeyNextStepCount)
		default:
			system.SetError(fmt.Sprintf("Anim %d is not implemented yet", a.AnimType))
			return
		}
	}
}
