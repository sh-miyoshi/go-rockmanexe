package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawVulcan struct {
	img []int
}

func (p *DrawVulcan) Init() error {
	path := common.ImagePath + "battle/skill/"

	fname := path + "バルカン.png"
	p.img = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 66, 50, p.img); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (p *DrawVulcan) End() {
	for i := 0; i < len(p.img); i++ {
		dxlib.DeleteGraph(p.img[i])
	}
	p.img = []int{}
}

func (p *DrawVulcan) Draw(viewPos common.Point, count int) {
	imgNo := 0
	if count > resources.SkillVulcanDelay*1 {
		imgNo = (count/(resources.SkillVulcanDelay*5))%2 + 1
	}

	// Show body
	dxlib.DrawRotaGraph(viewPos.X+50, viewPos.Y-18, 1, 0, p.img[imgNo], true)
	// Show attack
	if imgNo != 0 {
		if imgNo%2 == 0 {
			dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y-10, 1, 0, p.img[3], true)
		} else {
			dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y-15, 1, 0, p.img[3], true)
		}
	}
}
