package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
)

type countBomb struct {
	ID  string
	Arg Argument

	count int
}

func newCountBomb(objID string, arg Argument) *countBomb {
	return &countBomb{
		ID:  objID,
		Arg: arg,
	}
}

func (p *countBomb) Draw() {
	draw.String(100, 100, 0xff0000, "TODO: 落ちてくるアニメーション")
}

func (p *countBomb) Process() (bool, error) {
	const endCount = 60

	if p.count == 0 {
		field.SetBlackoutCount(endCount)
		SetChipNameDraw("カウントボム", true)
	}

	if p.count == endCount {
		pos := localanim.ObjAnimGetObjPos(p.Arg.OwnerID)
		obj := &object.CountBomb{}
		pm := object.ObjectParam{
			Pos:           common.Point{X: pos.X + 1, Y: pos.Y},
			HP:            50,
			OnwerCharType: objanim.ObjTypePlayer,
			Power:         int(p.Arg.Power),
		}
		if err := obj.Init(p.ID, pm); err != nil {
			return false, fmt.Errorf("count bomb create failed: %w", err)
		}
		atkID := localanim.ObjAnimNew(obj)
		localanim.ObjAnimAddActiveAnim(atkID)
		return true, nil
	}

	p.count++
	return false, nil
}

func (p *countBomb) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *countBomb) StopByOwner() {
}
