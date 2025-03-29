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

type AirHockeyDrawParam struct {
	PrevPos point.Point
	NextPos point.Point
}

type airhockey struct {
	ID   string
	Arg  Argument
	Core *processor.AirHockey
}

func newAirHockey(arg Argument, core skillcore.SkillCore) *airhockey {
	return &airhockey{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.AirHockey),
	}
}

func (p *airhockey) Draw() {
	// nothing to do at router
}

func (p *airhockey) Update() (bool, error) {
	return p.Core.Update()
}

func (p *airhockey) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeAirHockey,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	prev, current, next := p.Core.GetPos()
	drawPm := AirHockeyDrawParam{
		PrevPos: prev,
		NextPos: next,
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		Pos:       current,
		ExtraInfo: info.Marshal(),
	}
}

func (p *airhockey) StopByOwner() {
}

func (p *AirHockeyDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *AirHockeyDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
