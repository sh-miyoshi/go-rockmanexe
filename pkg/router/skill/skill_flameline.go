package skill

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type FlameLineDrawParam struct {
	Pillars []processor.FlamePillerParam
	Delay   int
}

type flameLine struct {
	ID   string
	Arg  Argument
	Core *processor.FlamePillarManager
}

func newFlameLine(arg Argument, core skillcore.SkillCore) *flameLine {
	return &flameLine{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.FlamePillarManager),
	}
}

func (p *flameLine) Draw() {
	// nothing to do at router
}

func (p *flameLine) Process() (bool, error) {
	return p.Core.Process()
}

func (p *flameLine) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
		AnimType:      routeranim.TypeFlameLine,
	}
	drawPm := FlameLineDrawParam{
		Pillars: p.Core.GetPillars(),
		Delay:   p.Core.GetDelay(),
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *flameLine) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *flameLine) GetEndCount() int {
	return p.Core.GetEndCount()
}

func (p *FlameLineDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *FlameLineDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}