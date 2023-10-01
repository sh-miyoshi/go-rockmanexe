package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	delayRecover = 1
)

type DrawRecover struct {
	imgRecover []int
}

func (p *DrawRecover) Init() error {
	path := common.ImagePath + "battle/skill/"

	fname := path + "リカバリー.png"
	p.imgRecover = make([]int, 8)
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 144, p.imgRecover); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (p *DrawRecover) End() {
	for i := 0; i < len(p.imgRecover); i++ {
		dxlib.DeleteGraph(p.imgRecover[i])
	}
	p.imgRecover = []int{}
}

func (p *DrawRecover) Draw(viewPos common.Point, count int) {
	n := (count / delayRecover) % len(p.imgRecover)
	if n >= 0 {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y, 1, 0, p.imgRecover[n], true)
	}
}
