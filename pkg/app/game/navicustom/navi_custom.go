package navicustom

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/fade"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/naviparts"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	stateOpening int = iota
	stateMain
	stateRun
)

var (
	state    int
	count    int
	imgBack  int = -1
	imgBoard int = -1
	// playerInfo *player.Player // TODO: 更新のタイミングで使う
	unsetParts []player.NaviCustomParts
)

func Init(plyr *player.Player) error {
	// playerInfo = plyr
	unsetParts = []player.NaviCustomParts{}
	for _, p := range plyr.AllNaviCustomParts {
		if !p.IsSet {
			unsetParts = append(unsetParts, p)
		}
	}

	state = stateOpening
	count = 0
	fname := common.ImagePath + "naviCustom/back.png"
	imgBack = dxlib.LoadGraph(fname)
	if imgBack == -1 {
		return fmt.Errorf("failed to load back image")
	}
	fname = common.ImagePath + "naviCustom/board.png"
	imgBoard = dxlib.LoadGraph(fname)
	if imgBoard == -1 {
		return fmt.Errorf("failed to load board image")
	}
	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	dxlib.DeleteGraph(imgBoard)
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, false)
	dxlib.DrawGraph(10, 30, imgBoard, true)

	switch state {
	case stateOpening:
		// Nothing to do
	case stateMain:
		// 実際にパーツを置いたりする
		for i, p := range unsetParts {
			parts := naviparts.Get(p.ID)
			dxlib.DrawBox(300, i*30+40, 400, i*30+65, dxlib.GetColor(16, 80, 104), true)
			draw.String(305, i*30+42, 0xFFFFFF, "%s", parts.Name)
		}
	case stateRun:
		// RUN
	}
}

func Process() {
	switch state {
	case stateOpening:
		if count == 0 {
			fade.In(30)
		}

		if count > 30 {
			stateChange(stateMain)
		}
	case stateMain:
		// TODO
	case stateRun:
		// TODO
	}

	count++
}

func stateChange(next int) {
	logger.Info("Change navu cutom state from %d to %d", state, next)
	state = next
	count = 0
}
