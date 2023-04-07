package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	aquamanStateInit int = iota
	aquamanStateAppear
	aquamanStateCreatePipe
	aquamanStateAttack
)

type aquaman struct {
	ID  string
	Arg Argument

	count         int
	state         int
	imgCharStand  []int
	imgCharCreate []int
	pos           common.Point
	atkID         string
}

func newAquaman(objID string, arg Argument) (*aquaman, error) {
	res := &aquaman{
		ID:    objID,
		Arg:   arg,
		state: aquamanStateInit,
		pos:   objanim.GetObjPos(arg.OwnerID),
	}

	fname := common.ImagePath + "battle/character/アクアマン_stand.png"
	res.imgCharStand = make([]int, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, res.imgCharStand); res == -1 {
		return nil, fmt.Errorf("failed to load image: %s", fname)
	}

	fname = common.ImagePath + "battle/character/アクアマン_create.png"
	res.imgCharCreate = make([]int, 1)
	if res := dxlib.LoadDivGraph(fname, 1, 1, 1, 80, 92, res.imgCharCreate); res == -1 {
		return nil, fmt.Errorf("failed to load image: %s", fname)
	}

	return res, nil
}

func (p *aquaman) Draw() {
	view := battlecommon.ViewPos(p.pos)
	xflip := int32(dxlib.TRUE)

	switch p.state {
	case aquamanStateInit:
	case aquamanStateAppear:
		const delay = 8
		if p.count > 20 {
			imgNo := (p.count / delay) % len(p.imgCharStand)
			dxlib.DrawRotaGraph(view.X+35, view.Y, 1, 0, p.imgCharStand[imgNo], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
		}
	case aquamanStateCreatePipe:
		imgNo := p.count
		if imgNo >= len(p.imgCharCreate) {
			imgNo = len(p.imgCharCreate) - 1
		}
		dxlib.DrawRotaGraph(view.X+35, view.Y, 1, 0, p.imgCharCreate[imgNo], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
	case aquamanStateAttack:
		dxlib.DrawRotaGraph(view.X+35, view.Y, 1, 0, p.imgCharCreate[len(p.imgCharCreate)-1], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
	}
}

func (p *aquaman) Process() (bool, error) {
	switch p.state {
	case aquamanStateInit:
		field.SetBlackoutCount(300)
		setChipNameDraw("アクアマン")
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
				Pos:           common.Point{X: p.pos.X + 1, Y: p.pos.Y},
				HP:            500,
				OnwerCharType: objanim.ObjTypePlayer,
				AttackNum:     1,
				Interval:      50,
				Power:         int(p.Arg.Power),
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
			field.SetBlackoutCount(0)
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
		DrawType: anim.DrawTypeSkill,
	}
}

func (p *aquaman) StopByOwner() {
	// Nothing to do after throwing
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
