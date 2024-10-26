package skill

import (
	"bytes"
	"encoding/gob"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type ShockWaveDrawParam struct {
	Speed    int
	Direct   int
	InitWait int
}

type shockWave struct {
	ID   string
	Arg  Argument
	Core *processor.ShockWave
}

func newShockWave(arg Argument, core skillcore.SkillCore) *shockWave {
	return &shockWave{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.ShockWave),
	}
}

func (p *shockWave) Draw() {
	// nothing to do at router
}

func (p *shockWave) Update() (bool, error) {
	return p.Core.Update()
}

func (p *shockWave) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeShockWave,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}
	pm := p.Core.GetParam()
	drawPm := ShockWaveDrawParam{
		Speed:    pm.Speed,
		Direct:   pm.Direct,
		InitWait: pm.InitWait,
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Core.GetPos(),
		DrawType:  anim.DrawTypeSkill,
		ExtraInfo: info.Marshal(),
	}
}

func (p *shockWave) StopByOwner() {
}

func (p *ShockWaveDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *ShockWaveDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
