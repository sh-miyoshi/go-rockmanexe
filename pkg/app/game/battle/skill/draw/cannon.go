package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayCannonAtk  = 2
	delayCannonBody = 6
)

type DrawCannon struct {
	imgBody [resources.SkillTypeCannonMax][]int
	imgAtk  [resources.SkillTypeCannonMax][]int
}

func (p *DrawCannon) Init() error {
	path := common.ImagePath + "battle/skill/"

	tmp := make([]int, 24)
	fname := path + "キャノン_atk.png"
	if res := dxlib.LoadDivGraph(fname, 24, 8, 3, 120, 140, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 8; i++ {
		p.imgAtk[0] = append(p.imgAtk[0], tmp[i])
		p.imgAtk[1] = append(p.imgAtk[1], tmp[i+8])
		p.imgAtk[2] = append(p.imgAtk[2], tmp[i+16])
	}
	fname = path + "キャノン_body.png"
	if res := dxlib.LoadDivGraph(fname, 15, 5, 3, 46, 40, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 5; i++ {
		p.imgBody[0] = append(p.imgBody[0], tmp[i])
		p.imgBody[1] = append(p.imgBody[1], tmp[i+5])
		p.imgBody[2] = append(p.imgBody[2], tmp[i+10])
	}
	return nil
}

func (p *DrawCannon) End() {
	for i := 0; i < 3; i++ {
		for j := 0; j < len(p.imgAtk[i]); j++ {
			dxlib.DeleteGraph(p.imgAtk[i][j])
		}
		p.imgAtk[i] = []int{}
		for j := 0; j < len(p.imgBody[i]); j++ {
			dxlib.DeleteGraph(p.imgBody[i][j])
		}
		p.imgBody[i] = []int{}
	}
}

func (p *DrawCannon) Draw(cannonType int, viewPos common.Point, count int) {
	n := count / delayCannonBody
	if n < len(p.imgBody[cannonType]) {
		if n >= 3 {
			viewPos.X -= 15
		}

		dxlib.DrawRotaGraph(viewPos.X+48, viewPos.Y-12, 1, 0, p.imgBody[cannonType][n], true)
	}

	n = (count - 15) / delayCannonAtk
	if n >= 0 && n < len(p.imgAtk[cannonType]) {
		dxlib.DrawRotaGraph(viewPos.X+90, viewPos.Y-10, 1, 0, p.imgAtk[cannonType][n], true)
	}
}
