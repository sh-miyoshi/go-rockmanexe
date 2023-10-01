package draw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/net"
	drawskill "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/netbattle/draw/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type animDraw struct {
	drawCannonInst   skilldraw.DrawCannon
	drawMiniBombInst skilldraw.DrawMiniBomb
	drawRecover      skilldraw.DrawRecover
	drawShockWave    skilldraw.DrawShockWave
	drawSpreadGun    drawskill.DrawSpreadGun
	drawSpreadHit    drawskill.DrawSpreadHit
	drawSword        drawskill.DrawSword
	drawVulcan       drawskill.DrawVulcan
	drawWideShot     drawskill.DrawWideShot
}

func (d *animDraw) Init() error {
	if err := d.drawCannonInst.Init(); err != nil {
		return fmt.Errorf("draw cannon init failed: %w", err)
	}
	if err := d.drawMiniBombInst.Init(); err != nil {
		return fmt.Errorf("draw minibomb init failed: %w", err)
	}
	if err := d.drawRecover.Init(); err != nil {
		return fmt.Errorf("draw recover init failed: %w", err)
	}
	if err := d.drawShockWave.Init(); err != nil {
		return fmt.Errorf("draw shock wave init failed: %w", err)
	}
	if err := d.drawSpreadGun.Init(); err != nil {
		return fmt.Errorf("draw spread gun init failed: %w", err)
	}
	if err := d.drawSpreadHit.Init(); err != nil {
		return fmt.Errorf("draw spread hit init failed: %w", err)
	}
	if err := d.drawSword.Init(); err != nil {
		return fmt.Errorf("draw sword init failed: %w", err)
	}
	if err := d.drawVulcan.Init(); err != nil {
		return fmt.Errorf("draw vulcan init failed: %w", err)
	}
	if err := d.drawWideShot.Init(); err != nil {
		return fmt.Errorf("draw wide shot init failed: %w", err)
	}

	return nil
}

func (d *animDraw) End() {
	d.drawCannonInst.End()
	d.drawMiniBombInst.End()
	d.drawRecover.End()
	d.drawShockWave.End()
	d.drawSpreadGun.End()
	d.drawSpreadHit.End()
	d.drawSword.End()
	d.drawVulcan.End()
	d.drawWideShot.End()
}

func (d *animDraw) Draw() {
	ginfo := net.GetInst().GetGameInfo()
	for _, a := range ginfo.Anims {
		pos := battlecommon.ViewPos(a.Pos)

		switch a.AnimType {
		case anim.TypeCannonNormal:
			d.drawCannonInst.Draw(resources.SkillTypeNormalCannon, pos, a.ActCount)
		case anim.TypeCannonHigh:
			d.drawCannonInst.Draw(resources.SkillTypeHighCannon, pos, a.ActCount)
		case anim.TypeCannonMega:
			d.drawCannonInst.Draw(resources.SkillTypeMegaCannon, pos, a.ActCount)
		case anim.TypeMiniBomb:
			target := common.Point{X: a.Pos.X + 3, Y: a.Pos.Y}
			d.drawMiniBombInst.Draw(a.Pos, target, a.ActCount)
		case anim.TypeRecover:
			d.drawRecover.Draw(pos, a.ActCount)
		case anim.TypeShockWave:
			d.drawShockWave.Draw(pos, a.ActCount, 3, common.DirectRight) // debug
		case anim.TypeSpreadGun:
			d.drawSpreadGun.Draw(pos, a.ActCount)
		case anim.TypeSpreadHit:
			d.drawSpreadHit.Draw(pos, a.ActCount)
		case anim.TypeSword:
			d.drawSword.Draw(skill.TypeSword, pos, a.ActCount)
		case anim.TypeWideSword:
			d.drawSword.Draw(skill.TypeWideSword, pos, a.ActCount)
		case anim.TypeLongSword:
			d.drawSword.Draw(skill.TypeLongSword, pos, a.ActCount)
		case anim.TypeVulcan:
			d.drawVulcan.Draw(pos, a.ActCount)
		case anim.TypeWideShot:
			d.drawWideShot.Draw(a.Pos, a.ActCount)
		default:
			common.SetError(fmt.Sprintf("Anim %d is not implemented yet", a.AnimType))
			return
		}
	}
}
