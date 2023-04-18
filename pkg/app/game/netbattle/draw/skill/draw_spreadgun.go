package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawSpreadGun struct {
	imgBody []int
	imgAtk  []int
}

type DrawSpreadHit struct {
	img []int
}

func (p *DrawSpreadGun) Init() error {
	path := common.ImagePath + "battle/skill/"

	fname := path + "スプレッドガン_atk.png"
	p.imgAtk = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 75, 76, p.imgAtk); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = path + "スプレッドガン_body.png"
	p.imgBody = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 56, 76, p.imgBody); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (p *DrawSpreadGun) End() {
	for j := 0; j < len(p.imgAtk); j++ {
		dxlib.DeleteGraph(p.imgAtk[j])
	}
	p.imgAtk = []int{}
	for j := 0; j < len(p.imgBody); j++ {
		dxlib.DeleteGraph(p.imgBody[j])
	}
	p.imgBody = []int{}
}

func (p *DrawSpreadGun) Draw(viewPos common.Point, count int) {
	// TODO: 定義場所を統一する
	const delaySpreadGun = 3

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

func (p *DrawSpreadHit) End() {
	for i := 0; i < len(p.img); i++ {
		dxlib.DeleteGraph(p.img[i])
	}
	p.img = []int{}
}

func (p *DrawSpreadHit) Draw(viewPos common.Point, count int) {
	if count < len(p.img) {
		dxlib.DrawRotaGraph(viewPos.X, viewPos.Y+15, 1, 0, p.img[count], true)
	}
}
