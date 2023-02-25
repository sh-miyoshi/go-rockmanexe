package skill

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
)

const (
	TypeNormalCannon int = iota
	TypeHighCannon
	TypeMegaCannon

	TypeCannonMax
)

const (
	delayCannonAtk   = 2
	delayCannonBody  = 6
	imgCannonBodyNum = 5
	imgCannonAtkNum  = 8
)

type cannon struct {
	ID   string
	Type int
	Arg  Argument

	count int
}

func newCannon(objID string, cannonType int, arg Argument) *cannon {
	return &cannon{
		ID:   objID,
		Type: cannonType,
		Arg:  arg,
	}
}

func (p *cannon) Draw() {
	// nothing to do at router
}

func (p *cannon) Process() (bool, error) {
	p.count++

	if p.count == 20 {
		// TODO add damage
	}

	max := imgCannonBodyNum * delayCannonBody
	if max < imgCannonAtkNum*delayCannonAtk+15 {
		max = imgCannonAtkNum*delayCannonAtk + 15
	}

	if p.count > max {
		return true, nil
	}
	return false, nil
}

func (p *cannon) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *cannon) StopByOwner() {
	anim.Delete(p.ID)
}
