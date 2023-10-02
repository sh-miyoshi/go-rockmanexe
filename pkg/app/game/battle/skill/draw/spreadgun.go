package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const delaySpreadGun = 3

type DrawSpreadGun struct {
	imgBody []int
	imgAtk  []int
}

type DrawSpreadHit struct {
	img []int
}

func (p *DrawSpreadGun) Init() {
	p.imgAtk = imgSpreadGunAtk
	p.imgBody = imgSpreadGunBody
}

func (p *DrawSpreadGun) Draw(viewPos common.Point, count int) {
	n := count / delaySpreadGun

	// Show body
	if n < len(p.imgBody) {
		dxlib.DrawRotaGraph(viewPos.X+50, viewPos.Y-18, 1, 0, p.imgBody[n], true)
	}

	// Show atk
	n = (count - 4) / delaySpreadGun
	if n >= 0 && n < len(p.imgAtk) {
		dxlib.DrawRotaGraph(viewPos.X+100, viewPos.Y-20, 1, 0, p.imgAtk[n], true)
	}
}

func (p *DrawSpreadHit) Init() error {
	p.img = make([]int, 6)
	fname := common.ImagePath + "battle/effect/spread_and_bamboo_hit.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, p.img); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (p *DrawSpreadHit) Draw(viewPos common.Point, count int) {
	if count < len(p.img) {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y+15, 1, 0, p.img[count], true)
	}
}
