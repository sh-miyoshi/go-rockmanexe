package draw

import (
	"fmt"

	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/skill"
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
}

func (d *animDraw) Init() error {
	if err := skilldraw.LoadImages(); err != nil {
		return fmt.Errorf("failed to load skill image: %w", err)
	}

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
			d.drawMiniBombInst.Draw(a.Pos, drawPm.Target, a.ActCount, drawPm.EndCount)
		case anim.TypeRecover:
			d.drawRecover.Draw(pos, a.ActCount)
		case anim.TypeShockWave:
			var drawPm skill.ShockWaveDrawParam
			drawPm.Unmarshal(a.DrawParam)
			d.drawShockWave.Draw(pos, a.ActCount, drawPm.Speed, drawPm.Direct)
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
		default:
			system.SetError(fmt.Sprintf("Anim %d is not implemented yet", a.AnimType))
			return
		}
	}
}
