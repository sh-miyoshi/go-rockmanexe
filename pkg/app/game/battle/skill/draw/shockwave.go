package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawShockWave struct {
	imgShockWave []int
}

func (p *DrawShockWave) Init() error {
	path := common.ImagePath + "battle/skill/"

	fname := path + "ショックウェーブ.png"
	p.imgShockWave = make([]int, 7)
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 100, 140, p.imgShockWave); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (p *DrawShockWave) End() {
	for i := 0; i < len(p.imgShockWave); i++ {
		dxlib.DeleteGraph(p.imgShockWave[i])
	}
	p.imgShockWave = []int{}
}

func (p *DrawShockWave) Draw(viewPos common.Point, count int, speed int, direct int) {
	n := (count / speed) % len(p.imgShockWave)
	if direct == common.DirectLeft {
		flag := int32(dxlib.TRUE)
		dxopts := dxlib.DrawRotaGraphOption{
			ReverseXFlag: &flag,
		}
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, p.imgShockWave[n], true, dxopts)
	} else if direct == common.DirectRight {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, p.imgShockWave[n], true)
	}
}
