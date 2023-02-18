package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type bambooLance struct {
	ID  string
	Arg Argument

	count    int
	imgSizeX int
}

func newBambooLance(objID string, arg Argument) *bambooLance {
	var sx, sy int
	dxlib.GetGraphSize(imgBambooLance[0], &sx, &sy)

	return &bambooLance{
		ID:       objID,
		Arg:      arg,
		imgSizeX: sx,
	}
}

func (p *bambooLance) Draw() {
	xreverse := int32(dxlib.TRUE)
	opt := dxlib.DrawRotaGraphOption{
		ReverseXFlag: &xreverse,
	}

	xd := p.count * 25
	if xd > field.PanelSize.X {
		xd = field.PanelSize.X
	}
	x := common.ScreenSize.X + p.imgSizeX/2 - xd
	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		v := battlecommon.ViewPos(common.Point{X: 0, Y: y})
		dxlib.DrawRotaGraph(x, v.Y+field.PanelSize.Y/2, 1, 0, imgBambooLance[0], true, opt)
	}
}

func (p *bambooLance) Process() (bool, error) {
	p.count++

	if p.count == 5 {
		dm := damage.Damage{
			Pos:           common.Point{X: battlecommon.FieldNum.X - 1},
			Power:         int(p.Arg.Power),
			TTL:           5,
			TargetType:    p.Arg.TargetType,
			HitEffectType: effect.TypeBambooHit,
			ShowHitArea:   false,
			BigDamage:     true,
			PushLeft:      1,
			DamageType:    damage.TypeWood,
		}
		for y := 0; y < battlecommon.FieldNum.Y; y++ {
			dm.Pos.Y = y
			damage.New(dm)
		}
	}

	if p.count > 10 {
		return true, nil
	}

	return false, nil
}

func (p *bambooLance) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *bambooLance) StopByOwner() {
	// Nothing to do after throwing
}
