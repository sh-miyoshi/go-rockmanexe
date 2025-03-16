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

type BoomerangDrawParam struct {
	PrevPos       point.Point
	NextPos       point.Point
	NextStepCount int
}

type boomerang struct {
	ID   string
	Arg  Argument
	Core *processor.Boomerang
}

func newBoomerang(arg Argument, core skillcore.SkillCore) *boomerang {
	return &boomerang{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.Boomerang),
	}
}

func (p *boomerang) Draw() {
	// nothing to do at router
}

func (p *boomerang) Update() (bool, error) {
	return p.Core.Update()
}

func (p *boomerang) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
		AnimType:      routeranim.TypeBoomerang,
	}

	prev, current, next := p.Core.GetPos()
	nextStepCnt := p.Core.GetNextStepCount()

	drawPm := BoomerangDrawParam{
		PrevPos:       prev,
		NextPos:       next,
		NextStepCount: nextStepCnt,
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		Pos:       current,
		ExtraInfo: info.Marshal(),
	}
}

func (p *boomerang) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}

func (p *BoomerangDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *BoomerangDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
