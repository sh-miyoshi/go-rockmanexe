package skill

import (
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore/processor"
	routeranim "github.com/sh-miyoshi/go-rockmanexe/pkg/router/anim"
)

type spreadGun struct {
	ID   string
	Arg  Argument
	Core *processor.SpreadGun
}

type spreadHit struct {
	ID   string
	Arg  Argument
	Core processor.SpreadHit
}

func newSpreadGun(arg Argument, core skillcore.SkillCore) *spreadGun {
	return &spreadGun{
		ID:   arg.AnimObjID,
		Arg:  arg,
		Core: core.(*processor.SpreadGun),
	}
}

func (p *spreadGun) Draw() {
	// nothing to do at router
}

func (p *spreadGun) Update() (bool, error) {
	res, err := p.Core.Update()
	if err != nil {
		return false, err
	}
	for _, hit := range p.Core.PopSpreadHits() {
		p.Arg.Manager.AnimNew(&spreadHit{
			ID:   uuid.New().String(),
			Arg:  p.Arg,
			Core: hit,
		})
	}

	return res, nil
}

func (p *spreadGun) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeSpreadGun,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Arg.Manager.ObjAnimGetObjPos(p.Arg.OwnerObjectID),
		DrawType:  anim.DrawTypeEffect,
		ExtraInfo: info.Marshal(),
	}
}

func (p *spreadGun) StopByOwner() {
	if p.Core.GetCount() < 5 {
		p.Arg.Manager.AnimDelete(p.ID)
	}
}

func (p *spreadHit) Draw() {
	// nothing to do at router
}

func (p *spreadHit) Update() (bool, error) {
	return p.Core.Update()
}

func (p *spreadHit) GetParam() anim.Param {
	info := routeranim.NetInfo{
		AnimType:      routeranim.TypeSpreadHit,
		OwnerClientID: p.Arg.OwnerClientID,
		ActCount:      p.Core.GetCount(),
	}

	return anim.Param{
		ObjID:     p.ID,
		Pos:       p.Core.Pos,
		DrawType:  anim.DrawTypeEffect,
		ExtraInfo: info.Marshal(),
	}
}
