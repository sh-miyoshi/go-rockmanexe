package skilldraw

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawAquaman struct {
	imgCharStand  []int
	imgCharCreate []int
}

func (p *DrawAquaman) Init() {
	p.imgCharStand = imgAquamanCharStand
	p.imgCharCreate = imgAquamanCharCreate
}

func (p *DrawAquaman) Draw(viewPos common.Point, count int, state int) {
	xflip := int32(dxlib.TRUE)

	switch state {
	case resources.SkillAquamanStateInit:
	case resources.SkillAquamanStateAppear:
		const delay = 8
		if count > 20 {
			imgNo := (count / delay) % len(p.imgCharStand)
			dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y, 1, 0, p.imgCharStand[imgNo], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
		}
	case resources.SkillAquamanStateCreatePipe:
		imgNo := count
		if imgNo >= len(p.imgCharCreate) {
			imgNo = len(p.imgCharCreate) - 1
		}
		dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y, 1, 0, p.imgCharCreate[imgNo], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
	case resources.SkillAquamanStateAttack:
		dxlib.DrawRotaGraph(viewPos.X+35, viewPos.Y, 1, 0, p.imgCharCreate[len(p.imgCharCreate)-1], true, dxlib.DrawRotaGraphOption{ReverseXFlag: &xflip})
	}
}
