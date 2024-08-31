package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
)

type skillHellsRolling struct {
	ID  string
	Arg skillcore.Argument

	count int
}

func newForteHellsRolling(objID string, arg skillcore.Argument, _ skillcore.SkillCore) *skillHellsRolling {
	return &skillHellsRolling{
		ID:  objID,
		Arg: arg,
		// TODO: Core
	}
}

func (p *skillHellsRolling) Draw() {
	// p.drawer.Draw()
}

func (p *skillHellsRolling) Process() (bool, error) {
	p.count++

	return false, nil
}

func (p *skillHellsRolling) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *skillHellsRolling) StopByOwner() {
	localanim.AnimDelete(p.ID)
}
