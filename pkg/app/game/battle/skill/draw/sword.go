package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const delaySword = 3

type DrawSword struct {
	img [resources.SkillTypeSwordMax][]int
}

func (p *DrawSword) Init() error {
	path := common.ImagePath + "battle/skill/"

	fname := path + "ソード.png"
	tmp := make([]int, 12)
	if res := dxlib.LoadDivGraph(fname, 12, 4, 3, 160, 150, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 0; i < 4; i++ {
		// Note: In the image, the order of wide sword and long sword is swapped.
		p.img[0] = append(p.img[0], tmp[i])
		p.img[1] = append(p.img[1], tmp[i+8])
		p.img[2] = append(p.img[2], tmp[i+4])
	}

	return nil
}

func (p *DrawSword) End() {
	for i := 0; i < 3; i++ {
		for j := 0; j < len(p.img[i]); j++ {
			dxlib.DeleteGraph(p.img[i][j])
		}
		p.img[i] = []int{}
	}
}

func (p *DrawSword) Draw(swordType int, viewPos common.Point, count int) {
	n := (count - 5) / delaySword
	if n >= 0 && n < len(p.img[swordType]) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y, 1, 0, p.img[swordType][n], true)
	}
}
