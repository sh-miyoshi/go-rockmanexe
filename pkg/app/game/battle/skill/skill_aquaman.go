package skill

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	skilldraw "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
)

type aquaman struct {
	ID   string
	Arg  skillcore.Argument
	Core *processor.Aquaman

	atkID  string
	drawer skilldraw.DrawAquaman
}

func newAquaman(objID string, arg skillcore.Argument, core skillcore.SkillCore) *aquaman {
	return &aquaman{
		ID:   objID,
		Arg:  arg,
		Core: core.(*processor.Aquaman),
	}
}

func (p *aquaman) Draw() {
	view := battlecommon.ViewPos(p.Core.GetPos())
	p.drawer.Draw(view, p.Core.GetCount(), p.Core.GetState())
}

func (p *aquaman) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}

	pipeObj := p.Core.PopWaterPipe()
	if pipeObj != nil {
		obj := &object.WaterPipe{}
		pm := object.ObjectParam{
			Pos:           pipeObj.Pos,
			HP:            500,
			OnwerCharType: objanim.ObjTypePlayer,
			AttackNum:     1,
			Interval:      50,
			Power:         int(p.Arg.Power),
		}
		if err := obj.Init(p.ID, pm); err != nil {
			return false, errors.Wrap(err, "water pipe create failed")
		}
		p.atkID = localanim.ObjAnimNew(obj)
		localanim.ObjAnimAddActiveAnim(p.atkID)
	}

	if p.atkID != "" {
		if !localanim.AnimIsProcessing(p.atkID) {
			field.SetBlackoutCount(0)
			return true, nil
		}
	}

	return res, nil
}

func (p *aquaman) GetParam() anim.Param {
	return anim.Param{
		ObjID: p.ID,
	}
}

func (p *aquaman) StopByOwner() {
	// Nothing to do after throwing
}
