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
	imgBlocks      [3]int
	playerInfo     *player.Player
	unsetParts     []player.NaviCustomParts
	setParts       []player.NaviCustomParts
	itemList       list.ItemList
	selected       int
	setPointerPos  common.Point
)

func Init(plyr *player.Player) error {
	playerInfo = plyr
	initParts()

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
	fname = common.ImagePath + "naviCustom/block_white.png"
	imgBlocks[colorBlock(naviparts.ColorWhite)] = dxlib.LoadGraph(fname)
	fname = common.ImagePath + "naviCustom/block_yellow.png"
	imgBlocks[colorBlock(naviparts.ColorYellow)] = dxlib.LoadGraph(fname)
	fname = common.ImagePath + "naviCustom/block_pink.png"
	imgBlocks[colorBlock(naviparts.ColorPink)] = dxlib.LoadGraph(fname)
	for i, b := range imgBlocks {
		if b == -1 {
			return fmt.Errorf("failed to load block %d image", i)
		}
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgBack)
	dxlib.DeleteGraph(imgBoard)
	dxlib.DeleteGraph(imgListPointer)
	dxlib.DeleteGraph(imgSetPointer)
	for b := range imgBlocks {
		dxlib.DeleteGraph(b)
	}
}

func Draw() {
	dxlib.DrawGraph(0, 0, imgBack, false)
	dxlib.DrawGraph(10, 30, imgBoard, true)

	switch state {
	case stateOpening:
		// Nothing to do
	case stateMain:
		for i, name := range itemList.GetList() {
			drawPartsListItem(300, i*30+45, name)
		}
		dxlib.DrawGraph(280, itemList.GetPointer()*30+50, imgListPointer, true)

		// セット済みのパーツを描画
		for _, s := range setParts {
			parts := naviparts.Get(s.ID)
			drawBoardParts(common.Point{X: s.X, Y: s.Y}, parts)
		}

		// TODO: ミニウィンドウ

		if selected != -1 {
			parts := naviparts.Get(unsetParts[selected].ID)
			drawBoardParts(setPointerPos, parts)
			if (count/10)%3 != 0 {
				baseX := setPointerPos.X*40 + 34
				baseY := setPointerPos.Y*40 + 65
				dxlib.DrawGraph(baseX, baseY, imgSetPointer, true)
			}
		}

		// コマンドライン
		dxlib.DrawBox(32, 158, 240, 161, 0x282828, true)
		dxlib.DrawBox(32, 174, 240, 177, 0x282828, true)
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
				sound.On(resources.SECancel)
				selected = -1
				return
			}

			if inputs.CheckKey(inputs.KeyEnter) == 1 {
				if checkSet() {
					sound.On(resources.SEMenuEnter)
					newInfo := unsetParts[selected]
					newInfo.IsSet = true
					newInfo.X = setPointerPos.X
					newInfo.Y = setPointerPos.Y

					playerInfo.UpdateNaviCustomParts(unsetParts[selected].ObjID, newInfo)
					initParts()
					selected = -1
				} else {
					sound.On(resources.SEDenied)
				}
				return
			}

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
	dxlib.DrawBox(x-2, y-1, x+122, y+26, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(x, y, x+120, y+25, dxlib.GetColor(16, 80, 104), true)
	draw.String(x+5, y+2, 0xFFFFFF, "%s", name)
}

func colorBlock(color int) int {
	switch color {
	case naviparts.ColorWhite:
		return 0
	case naviparts.ColorYellow:
		return 1
	case naviparts.ColorPink:
		return 2
	}

	common.SetError(fmt.Sprintf("カラーコード %d に対するブロックは存在しません", color))
	return 0
}

func drawBoardParts(basePos common.Point, parts naviparts.NaviParts) {
	baseX := basePos.X*40 + 34
	baseY := basePos.Y*40 + 65

	for _, b := range parts.Blocks {
		x := basePos.X + b.X
		y := basePos.Y + b.Y
		if x < 0 || x >= boardSize || y < 0 || y >= boardSize {
			continue
		}

		if parts.IsPlusParts {
			dxlib.DrawGraph(b.X*40+baseX+4, b.Y*40+baseY+4, imgBlocks[colorBlock(parts.Color)], true)
		} else {
			dxlib.DrawBox(b.X*40+baseX+4, b.Y*40+baseY+4, (b.X+1)*40+baseX, (b.Y+1)*40+baseY, naviparts.GetColorCode(parts.Color), true)
		}
	}
}

func checkSet() bool {
	parts := naviparts.Get(unsetParts[selected].ID)

	// セットするパートがボード外にはみ出ていないか
	for _, b := range parts.Blocks {
		x := setPointerPos.X + b.X
		y := setPointerPos.Y + b.Y
		if x < 0 || x >= boardSize || y < 0 || y >= boardSize {
			return false
		}
	}

	// パーツ同士が重なっていないか
	// TODO

	return true // セットできる
}

func initParts() {
	unsetParts = []player.NaviCustomParts{}
	setParts = []player.NaviCustomParts{}
	names := []string{}
	for _, p := range playerInfo.AllNaviCustomParts {
		if !p.IsSet {
			unsetParts = append(unsetParts, p)
			parts := naviparts.Get(p.ID)
			names = append(names, parts.Name)
		} else {
			setParts = append(setParts, p)
		}
	}
	itemList.SetList(names, maxListNum)
}
