package skill

import (
	"bytes"
	"encoding/gob"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/router/gameinfo"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/queue"
)

type VulcanDrawParam struct {
	Delay int
}

type vulcan struct {
	ID   string
	Arg  Argument
	Core (*processor.Vulcan)
}

func newVulcan(arg Argument, core skillcore.SkillCore) *vulcan {
	return &vulcan{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.Vulcan),
	}
}

func (p *vulcan) Draw() {
	// nothing to do at router
}

func (p *vulcan) Process() (bool, error) {
	res, err := p.Core.Process()
	if err != nil {
		return false, err
	}
	for _, eff := range p.Core.PopEffects() {
		queue.Push(p.Arg.QueueIDs[gameinfo.QueueTypeEffect], &gameinfo.Effect{
			ID:            uuid.New().String(),
			OwnerClientID: p.Arg.GameInfo.ClientID,
			Pos:           eff.Pos,
			Type:          eff.Type,
			RandRange:     eff.RandRange,
		})
	}

	return res, nil
}

func (p *vulcan) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		AnimType:      routeranim.TypeVulcan,
		ActCount:      p.Core.GetCount(),
	}
	drawPm := VulcanDrawParam{Delay: p.Core.GetDelay()}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       routeranim.ObjAnimGetObjPos(p.Arg.OwnerClientID, p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *vulcan) StopByOwner() {
	routeranim.AnimDelete(p.Arg.OwnerClientID, p.ID)
}

func (p *vulcan) GetEndCount() int {
	return p.Core.GetEndCount()
}

func (p *VulcanDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *VulcanDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
