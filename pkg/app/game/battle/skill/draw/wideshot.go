package skilldraw

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type DrawWideShot struct {
	imgBody  []int
	imgBegin []int
	imgMove  []int
}

func (p *DrawWideShot) Init() error {
	path := common.ImagePath + "battle/skill/"

	fname := path + "ワイドショット_body.png"
	p.imgBody = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 56, 66, p.imgBody); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = path + "ワイドショット_begin.png"
	p.imgBegin = make([]int, 4)
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 90, 147, p.imgBegin); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	fname = path + "ワイドショット_move.png"
	p.imgMove = make([]int, 3)
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 90, 148, p.imgMove); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	return nil
}

func (p *DrawWideShot) End() {
	for i := 0; i < len(p.imgBody); i++ {
		dxlib.DeleteGraph(p.imgBody[i])
	}
	p.imgBody = []int{}

	for i := 0; i < len(p.imgBegin); i++ {
		dxlib.DeleteGraph(p.imgBegin[i])
	}
	p.imgBegin = []int{}

	for i := 0; i < len(p.imgMove); i++ {
		dxlib.DeleteGraph(p.imgMove[i])
	}
	p.imgMove = []int{}
}

func (p *DrawWideShot) Draw(pos common.Point, count int, direct int, showBody bool, nextStepCount int, state int) {
	opt := dxlib.DrawRotaGraphOption{}
	ofs := 1
	if direct == common.DirectLeft {
		xflip := int32(dxlib.TRUE)
		opt.ReverseXFlag = &xflip
		ofs = -1
	}
	view := battlecommon.ViewPos(pos)

	switch state {
	case resources.SkillWideShotStateBegin:
		n := (count / resources.SkillWideShotDelay)

		if n < len(p.imgBody) && showBody {
			dxlib.DrawRotaGraph(view.X+40, view.Y-13, 1, 0, p.imgBody[n], true, opt)
		}
		if n >= len(p.imgBegin) {
			n = len(p.imgBegin) - 1
		}
		dxlib.DrawRotaGraph(view.X+62*ofs, view.Y+20, 1, 0, p.imgBegin[n], true, opt)
	case resources.SkillWideShotStateMove:
		n := (count / resources.SkillWideShotDelay) % len(p.imgMove)
		next := pos.X + 1
		prev := pos.X - 1
		if direct == common.DirectLeft {
			next, prev = prev, next
		}

		c := count % nextStepCount
		if c != 0 {
			ofsx := battlecommon.GetOffset(next, pos.X, prev, c, nextStepCount, battlecommon.PanelSize.X)
			dxlib.DrawRotaGraph(view.X+ofsx, view.Y+20, 1, 0, p.imgMove[n], true, opt)
		}
	}
}
