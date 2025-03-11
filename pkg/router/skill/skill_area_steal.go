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

type AreaStealDrawParam struct {
	State   int
	Targets []point.Point
}

type areaSteal struct {
	ID   string
	Arg  Argument
	Core *processor.AreaSteal
}

func newAreaSteal(arg Argument, core skillcore.SkillCore) *areaSteal {
	return &areaSteal{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.AreaSteal),
	}
}

func (p *areaSteal) Draw() {
	// nothing to do at router
}

func (p *areaSteal) Update() (bool, error) {
	return p.Core.Update()
}

func (p *areaSteal) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeAreaSteal,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}
	drawPm := AreaStealDrawParam{
		State:   p.Core.GetState(),
		Targets: p.Core.GetTargets(),
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID: p.ID,
		Pos:   p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),

		ExtraInfo: info.Marshal(),
	}
}

func (p *areaSteal) StopByOwner() {
}

func (p *AreaStealDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *AreaStealDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
