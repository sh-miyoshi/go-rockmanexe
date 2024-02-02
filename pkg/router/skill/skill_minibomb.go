package skill

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

type MiniBombDrawParam struct {
	EndCount int
	Target   point.Point
}

type miniBomb struct {
	ID   string
	Arg  Argument
	Core *processor.MiniBomb
}

func newMiniBomb(arg Argument, core skillcore.SkillCore) *miniBomb {
	return &miniBomb{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.MiniBomb),
	}
}

func (p *miniBomb) Draw() {
	// nothing to do at router
}

func (p *miniBomb) Process() (bool, error) {
	return p.Core.Process()
}

func (p *miniBomb) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeMiniBomb,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}
	current, target := p.Core.GetPointParams()
	drawPm := MiniBombDrawParam{
		EndCount: p.Core.GetEndCount(),
		Target:   target,
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		Pos:       current,
		DrawType:  anim.DrawTypeSkill,
		ExtraInfo: info.Marshal(),
	}
}

func (p *miniBomb) StopByOwner() {
	if p.Core.GetCount() < 5 {
		routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
	}
}

func (p *miniBomb) GetEndCount() int {
	return p.Core.GetEndCount()
}

func (p *MiniBombDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *MiniBombDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
