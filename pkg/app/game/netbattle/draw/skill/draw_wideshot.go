package skill

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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

func (p *DrawWideShot) Draw(pos common.Point, count int) {
	// TODO: 定義場所を統一する
	const delayWideShot = 4
	const nextStepCount = 8
	const stateSeparateCount = 1000
	// 順番をRouter側と合わせておく
	const (
		wideShotStateBegin int = iota
		wideShotStateMove
	)

	viewPos := battlecommon.ViewPos(pos)

	state := count / stateSeparateCount

	switch state {
	case wideShotStateBegin:
		if count == 0 {
			return
		}

		logger.Debug("count: %d", count)

		n := (count / delayWideShot)

		if n < len(p.imgBody) {
			dxlib.DrawRotaGraph(viewPos.X+40, viewPos.Y-13, 1, 0, p.imgBody[n], true)
		}
		if n >= len(p.imgBegin) {
			n = len(p.imgBegin) - 1
		}
		dxlib.DrawRotaGraph(viewPos.X+62, viewPos.Y+20, 1, 0, p.imgBegin[n], true)
	case wideShotStateMove:
		count -= stateSeparateCount
		n := (count / delayWideShot) % len(p.imgMove)
		ofsx := battlecommon.PanelSize.X*count/nextStepCount + (battlecommon.PanelSize.X / 2)
		dxlib.DrawRotaGraph(viewPos.X+ofsx, viewPos.Y+20, 1, 0, p.imgMove[n], true)
	}
}
