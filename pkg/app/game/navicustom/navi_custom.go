package navicustom

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/fade"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/naviparts"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/list"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	maxListNum = 5
	boardSize  = 5
)

const (
	stateOpening int = iota
	stateMain
	stateRun
)

var (
	state          int
	count          int
	imgBack        int = -1
	imgBoard       int = -1
	imgListPointer int = -1
	imgSetPointer  int = -1
	// playerInfo *player.Player // TODO: 更新のタイミングで使う
	unsetParts    []player.NaviCustomParts
	itemList      list.ItemList
	selected      int
	setPointerPos common.Point
)

func Init(plyr *player.Player) error {
	// playerInfo = plyr
	unsetParts = []player.NaviCustomParts{}
	names := []string{}
	for _, p := range plyr.AllNaviCustomParts {
		if !p.IsSet {
			unsetParts = append(unsetParts, p)
			parts := naviparts.Get(p.ID)
			names = append(names, parts.Name)
		}
	}
	itemList.SetList(names, maxListNum)

	state = stateOpening
	count = 0
	selected = -1
	setPointerPos = common.Point{X: 2, Y: 2}

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
	fname = common.ImagePath + "naviCustom/pointer.png"
	imgListPointer = dxlib.LoadGraph(fname)
	if imgListPointer == -1 {
		return fmt.Errorf("failed to load list pointer image")
	}
	fname = common.ImagePath + "naviCustom/pointer2.png"
	imgSetPointer = dxlib.LoadGraph(fname)
	if imgListPointer == -1 {
		return fmt.Errorf("failed to load set pointer image")
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	dxlib.DeleteGraph(imgBoard)
	dxlib.DeleteGraph(imgListPointer)
	dxlib.DeleteGraph(imgSetPointer)
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, false)
	dxlib.DrawGraph(10, 30, imgBoard, true)

	switch state {
	case stateOpening:
		// Nothing to do
	case stateMain:
		// 実際にパーツを置いたりする
		for i, name := range itemList.GetList() {
			drawPartsListItem(300, i*30+45, name)
		}
		dxlib.DrawGraph(280, itemList.GetPointer()*30+50, imgListPointer, true)

		// TODO: ミニウィンドウ

		if selected != -1 {
			baseX := setPointerPos.X*40 + 34
			baseY := setPointerPos.Y*40 + 65
			dxlib.DrawGraph(baseX, baseY, imgSetPointer, true)
			parts := naviparts.Get(unsetParts[selected].ID)
			for _, b := range parts.Blocks {
				// TODO
				dxlib.DrawBox(baseX+b.X*40, baseY+b.Y*40, baseX+(b.X+1)*40, baseY+(b.Y+1)*40, 0xFFFFFF, true)
			}
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
		if selected == -1 {
			selected = itemList.Process()
			if selected != -1 {
				sound.On(resources.SEMenuEnter)
			}
		} else {
			if inputs.CheckKey(inputs.KeyCancel) == 1 {
				selected = -1
				return
			}

			// TODO(決定)

			if inputs.CheckKey(inputs.KeyUp) == 1 && setPointerPos.Y > 0 {
				sound.On(resources.SECursorMove)
				setPointerPos.Y--
				return
			}
			if inputs.CheckKey(inputs.KeyDown) == 1 && setPointerPos.Y < boardSize-1 {
				sound.On(resources.SECursorMove)
				setPointerPos.Y++
				return
			}
			if inputs.CheckKey(inputs.KeyLeft) == 1 && setPointerPos.X > 0 {
				sound.On(resources.SECursorMove)
				setPointerPos.X--
				return
			}
			if inputs.CheckKey(inputs.KeyRight) == 1 && setPointerPos.X < boardSize-1 {
				sound.On(resources.SECursorMove)
				setPointerPos.X++
				return
			}
		}
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

func drawPartsListItem(x, y int, name string) {
	dxlib.DrawBox(x-2, y-1, x+102, y+26, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(x, y, x+100, y+25, dxlib.GetColor(16, 80, 104), true)
	draw.String(x+5, y+2, 0xFFFFFF, "%s", name)
}
