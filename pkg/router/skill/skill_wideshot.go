package skill

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type WideShotDrawParam struct {
	State         int
	NextStepCount int
	Direct        int
}

type wideShot struct {
	ID   string
	Arg  Argument
	Core *processor.WideShot
}

func newWideShot(arg Argument, core skillcore.SkillCore) *wideShot {
	return &wideShot{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.WideShot),
	}
}

func (p *wideShot) Draw() {
	// nothing to do at router
}

func (p *wideShot) Process() (bool, error) {
	return p.Core.Process()
}

func (p *wideShot) GetParam() anim.Param {
	pm := p.Core.GetParam()
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		AnimType:      routeranim.TypeWideShot,
		ActCount:      p.Core.GetCount(),
	}
	drawPm := WideShotDrawParam{
		State:         pm.State,
		NextStepCount: pm.NextStepCount,
		Direct:        pm.Direct,
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       pm.Pos,
		ExtraInfo: info.Marshal(),
	}
}

func (p *wideShot) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}

func (p *wideShot) GetEndCount() int {
	return p.Core.GetEndCount()
}

func (p *WideShotDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *WideShotDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
