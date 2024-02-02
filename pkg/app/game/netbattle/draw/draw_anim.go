package draw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
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
		case anim.TypeCannonNormal:
			d.drawCannonInst.Draw(resources.SkillCannon, pos, a.ActCount) // TODO: 要調整
		case anim.TypeCannonHigh:
			d.drawCannonInst.Draw(resources.SkillHighCannon, pos, a.ActCount)
		case anim.TypeCannonMega:
			d.drawCannonInst.Draw(resources.SkillMegaCannon, pos, a.ActCount)
		case anim.TypeMiniBomb:
			target := point.Point{X: a.Pos.X + 3, Y: a.Pos.Y}
			endCount := 60 // TODO: 要調整
			d.drawMiniBombInst.Draw(a.Pos, target, a.ActCount, endCount)
		case anim.TypeRecover:
			d.drawRecover.Draw(pos, a.ActCount)
		case anim.TypeShockWave:
			d.drawShockWave.Draw(pos, a.ActCount, 3, config.DirectRight) // debug
		case anim.TypeSpreadGun:
			d.drawSpreadGun.Draw(pos, a.ActCount)
		case anim.TypeSpreadHit:
			d.drawSpreadHit.Draw(pos, a.ActCount)
		case anim.TypeSword:
			d.drawSword.Draw(0, pos, a.ActCount) // TODO: 要調整
		case anim.TypeWideSword:
			d.drawSword.Draw(1, pos, a.ActCount)
		case anim.TypeLongSword:
			d.drawSword.Draw(2, pos, a.ActCount)
		case anim.TypeVulcan:
			d.drawVulcan.Draw(pos, a.ActCount)
		case anim.TypeWideShot:
			// TODO: refactoring
			state := a.ActCount / 1000
			a.ActCount -= state * 1000
			d.drawWideShot.Draw(a.Pos, a.ActCount, config.DirectRight, true, resources.SkillWideShotPlayerNextStepCount, state)
		default:
			system.SetError(fmt.Sprintf("Anim %d is not implemented yet", a.AnimType))
			return
		}
	}
}
