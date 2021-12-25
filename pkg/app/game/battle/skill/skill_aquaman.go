package skill

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
)

const (
	aquamanStateInit int = iota
	aquamanStateAppear
	aquamanStateCreatePipe
	aquamanStateAttack
)

type aquaman struct {
	ID         string
	OwnerID    string
	Power      uint
	TargetType int

	count         int
	state         int
	imgCharStand  []int32
	imgCharCreate []int32
	x, y          int
	atkID         string
}

func newAquaman(objID string, arg Argument) (*aquaman, error) {
	res := &aquaman{
		ID:         objID,
		OwnerID:    arg.OwnerID,
		Power:      arg.Power,
		TargetType: arg.TargetType,
		state:      aquamanStateInit,
	}

	res.x, res.y = objanim.GetObjPos(arg.OwnerID)

	fname := common.ImagePath + "battle/character/アクアマン_stand.png"
	res.imgCharStand = make([]int32, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, res.imgCharStand); res == -1 {
		return nil, fmt.Errorf("failed to load image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/アクアマン_create.png"
	res.imgCharCreate = make([]int32, 1)
	if res := dxlib.LoadDivGraph(fname, 1, 1, 1, 80, 92, res.imgCharCreate); res == -1 {
		return nil, fmt.Errorf("failed to load image: %s", fname)
	}

	return res, nil
}

func (p *aquaman) Draw() {
	px, py := battlecommon.ViewPos(p.x, p.y)
	xflip := int32(dxlib.TRUE)

	switch p.state {
	case aquamanStateInit:
	case aquamanStateAppear:
		const delay = 8
		if p.count > 20 {
			imgNo := (p.count / delay) % len(p.imgCharStand)
			dxlib.DrawRotaGraph(px+35, py, 1, 0, p.imgCharStand[imgNo], dxlib.TRUE, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
		}
	case aquamanStateCreatePipe:
		imgNo := p.count
		if imgNo >= len(p.imgCharCreate) {
			imgNo = len(p.imgCharCreate) - 1
		}
		dxlib.DrawRotaGraph(px+35, py, 1, 0, p.imgCharCreate[imgNo], dxlib.TRUE, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
	case aquamanStateAttack:
		dxlib.DrawRotaGraph(px+35, py, 1, 0, p.imgCharCreate[len(p.imgCharCreate)-1], dxlib.TRUE, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
	}
}

func (p *aquaman) Process() (bool, error) {
	switch p.state {
	case aquamanStateInit:
		field.SetBlackoutCount(300)
		p.setState(aquamanStateAppear)
		return false, nil
	case aquamanStateAppear:
		if p.count == 70 {
			p.setState(aquamanStateCreatePipe)
			return false, nil
		}
	case aquamanStateCreatePipe:
		if p.count == 10 {
			obj := &object.WaterPipe{}
			pm := object.ObjectParam{
				PosX:          p.x + 1,
				PosY:          p.y,
				HP:            500,
				OnwerCharType: objanim.ObjTypePlayer,
				AttackNum:     1,
				Interval:      50,
				Power:         int(p.Power),
			}
			if err := obj.Init(p.ID, pm); err != nil {
				return false, fmt.Errorf("water pipe create failed: %w", err)
			}
			p.atkID = objanim.New(obj)
			objanim.AddActiveAnim(p.atkID)

			p.setState(aquamanStateAttack)
			return false, nil
		}
	case aquamanStateAttack:
		if !objanim.IsProcessing(p.atkID) {
			p.end()
			return true, nil
		}
	}

	p.count++
	return false, nil
}

func (p *aquaman) GetParam() anim.Param {
	return anim.Param{
		ObjID:    p.ID,
		AnimType: anim.AnimTypeSkill,
	}
}

func (p *aquaman) setState(nextState int) {
	p.count = 0
	p.state = nextState
}

func (p *aquaman) end() {
	for _, img := range p.imgCharStand {
		dxlib.DeleteGraph(img)
	}
	for _, img := range p.imgCharCreate {
		dxlib.DeleteGraph(img)
	}
}
