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

type SwordDrawParam struct {
	SkillID int
	Delay   int
}

type sword struct {
	ID      string
	SkillID int
	Arg     Argument
	Core    *processor.Sword
}

func newSword(skillID int, arg Argument, core skillcore.SkillCore) *sword {
	return &sword{
		ID:      arg.AnimObjID,
		SkillID: skillID,
		Arg:     arg,
		Core:    core.(*processor.Sword),
	}
}

func (p *sword) Draw() {
	// nothing to do at router
}

func (p *sword) Update() (bool, error) {
	return p.Core.Update()
}

func (p *sword) GetParam() anim.Param {
	info := routeranim.NetInfo{
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}
	switch p.SkillID {
	case resources.SkillSword:
		info.AnimType = routeranim.TypeSword
	case resources.SkillWideSword:
		info.AnimType = routeranim.TypeWideSword
	case resources.SkillLongSword:
		info.AnimType = routeranim.TypeLongSword
	}
	drawPm := SwordDrawParam{
		SkillID: p.Core.GetID(),
		Delay:   p.Core.GetDelay(),
	}
	info.DrawParam = drawPm.Marshal()

	return anim.Param{
		ObjID:     p.ID,
		DrawType:  anim.DrawTypeSkill,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		ExtraInfo: info.Marshal(),
	}
}

func (p *sword) StopByOwner() {
	p.Arg.Manager.AnimDelete(p.ID)
}

func (p *SwordDrawParam) Marshal() []byte {
	buf := bytes.NewBuffer(nil)
	gob.NewEncoder(buf).Encode(p)
	return buf.Bytes()
}

func (p *SwordDrawParam) Unmarshal(data []byte) {
	buf := bytes.NewBuffer(data)
	_ = gob.NewDecoder(buf).Decode(p)
}
