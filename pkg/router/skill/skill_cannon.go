package skill

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type CannonDrawParam struct {
	Type int
}

type cannon struct {
	ID   string
	Arg  Argument
	Core *processor.Cannon
}

func newCannon(arg Argument, core skillcore.SkillCore) *cannon {
	return &cannon{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.Cannon),
	}
}

func (p *cannon) Draw() {
	// nothing to do at router
}

func (p *cannon) Process() (bool, error) {
	return p.Core.Process()
}

func (p *cannon) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	switch p.Core.GetCannonType() {
	case resources.SkillCannon:
		info.AnimType = routeranim.TypeCannonNormal
	case resources.SkillHighCannon:
		info.AnimType = routeranim.TypeCannonHigh
	case resources.SkillMegaCannon:
		info.AnimType = routeranim.TypeCannonMega
	}
	drawPm := CannonDrawParam{Type: p.Core.GetCannonType()}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *cannon) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}

func (p *cannon) GetEndCount() int {
	return p.Core.GetEndCount()
}

func (p *CannonDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *CannonDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
