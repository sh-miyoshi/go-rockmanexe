package draw

import (
	"fmt"

	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
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
}

func (d *animDraw) Init() error {
	if err := skilldraw.LoadImages(); err != nil {
		return fmt.Errorf("failed to load skill image: %w", err)
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

		switch a.AnimType {
		case anim.TypeCannonNormal, anim.TypeCannonHigh, anim.TypeCannonMega:
			var drawPm skill.CannonDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawCannonInst.Draw(drawPm.Type, pos, a.ActCount)
		case anim.TypeMiniBomb:
			var drawPm skill.MiniBombDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawMiniBombInst.Draw(a.Pos, drawPm.Target, a.ActCount, drawPm.LandCount)
		case anim.TypeRecover:
			d.drawRecover.Draw(pos, a.ActCount)
		case anim.TypeShockWave:
			var drawPm skill.ShockWaveDrawParam
			drawPm.Unmarshal(a.DrawParam)
			if a.ActCount >= drawPm.InitWait {
				d.drawShockWave.Draw(pos, a.ActCount, drawPm.Speed, drawPm.Direct)
			}
		case anim.TypeSpreadGun:
			d.drawSpreadGun.Draw(pos, a.ActCount)
		case anim.TypeSpreadHit:
			d.drawSpreadHit.Draw(pos, a.ActCount)
		case anim.TypeSword, anim.TypeWideSword, anim.TypeLongSword:
			var drawPm skill.SwordDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawSword.Draw(drawPm.Type, pos, a.ActCount, drawPm.Delay)
		case anim.TypeVulcan:
			var drawPm skill.VulcanDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawVulcan.Draw(pos, a.ActCount, drawPm.Delay)
		case anim.TypeWideShot:
			var drawPm skill.WideShotDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawWideShot.Draw(a.Pos, a.ActCount, drawPm.Direct, true, drawPm.NextStepCount, drawPm.State)
		case anim.TypeHeatShot, anim.TypeHeatV, anim.TypeHeatSide:
			d.drawHeatShot.Draw(pos, a.ActCount)
		case anim.TypeFlameLine:
			var drawPm skill.FlameLineDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawFlameLine.Draw(pos, a.ActCount, true, drawPm.Pillars, drawPm.Delay)
		case anim.TypeTornado:
			// Note: DrawParamで渡すようにしてもいいが、targetの決定アルゴリズムが変わることはないのでここに直接書く
			targetPos := point.Point{X: a.Pos.X + 2, Y: a.Pos.Y}
			target := battlecommon.ViewPos(targetPos)
			d.drawTornado.Draw(pos, target, a.ActCount)
		case anim.TypeBoomerang:
			var drawPm skill.BoomerangDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawBoomerang.Draw(drawPm.PrevPos, a.Pos, drawPm.NextPos, a.ActCount, drawPm.NextStepCount)
		case anim.TypeBambooLance:
			d.drawBamboolance.Draw(a.ActCount)
		case anim.TypeCrack:
			// no animation
		default:
			system.SetError(fmt.Sprintf("Anim %d is not implemented yet", a.AnimType))
			return
		}
	}
}
