package navicustom

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/fade"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/ncparts"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/list"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	maxListNum = 3
	boardSize  = 5
	runName    = "RUN"
	lineY      = 3 - 1
)

const (
	stateOpening int = iota
	stateUnsetPartsSelect
	stateBoardPartsSelect
	stateDeployment
	stateRun
	stateRunEnd
)

type partsInfo struct {
	rawData player.NaviCustomParts
	objID   string
}

var (
	state          int
	beforeState    int
	count          int
	imgBack        int = -1
	imgBoard       int = -1
	imgListPointer int = -1
	imgSetPointer  int = -1
	imgBlocks      [3]int
	playerInfo     *player.Player
	allParts       []partsInfo
	unsetParts     []partsInfo
	setParts       []partsInfo
	itemList       list.ItemList
	selected       int
	setPointerPos  common.Point
)

func Init(plyr *player.Player) error {
	playerInfo = plyr
	state = stateOpening
	beforeState = stateOpening
	count = 0
	selected = -1
	setPointerPos = common.Point{X: 2, Y: 2}
	allParts = []partsInfo{}
	for _, p := range plyr.AllNaviCustomParts {
		allParts = append(allParts, partsInfo{
			rawData: p,
			objID:   uuid.NewString(),
		})
	}

	initParts()

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
	imgBlocks[colorBlock(ncparts.ColorWhite)] = dxlib.LoadGraph(fname)
	fname = common.ImagePath + "naviCustom/block_yellow.png"
	imgBlocks[colorBlock(ncparts.ColorYellow)] = dxlib.LoadGraph(fname)
	fname = common.ImagePath + "naviCustom/block_pink.png"
	imgBlocks[colorBlock(ncparts.ColorPink)] = dxlib.LoadGraph(fname)
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
	if state == stateOpening && count < 2 {
		return
	}

	dxlib.DrawGraph(0, 0, imgBack, false)
	dxlib.DrawGraph(10, 30, imgBoard, true)

	for i := 0; i < maxListNum; i++ {
		c := i + itemList.GetScroll()
		name := itemList.GetList()[c]

		x := 300
		y := i*30 + 45
		if name == runName {
			dxlib.DrawBox(x-2, y-1, x+122, y+26, dxlib.GetColor(168, 192, 216), true)
			dxlib.DrawBox(x, y, x+120, y+25, dxlib.GetColor(0, 194, 33), true)
			draw.String(x+25, y+2, 0xFFFFFF, "ＲＵＮ！")
		} else {
			drawPartsListItem(x, y, name)
		}
	}
	dxlib.DrawGraph(280, itemList.GetPointer()*30+50, imgListPointer, true)

	// セット済みのパーツを描画
	for _, s := range setParts {
		parts := ncparts.Get(s.rawData.ID)
		drawBoardParts(common.Point{X: s.rawData.X, Y: s.rawData.Y}, parts)
	}

	// TODO: ミニウィンドウ

	// 選択しているパーツの描画
	if selected >= 0 && selected < len(unsetParts) {
		parts := ncparts.Get(unsetParts[selected].rawData.ID)
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

	// Information Panel
	x := 300
	y := 170
	dxlib.DrawBox(x-2, y-1, x+122, y+106, dxlib.GetColor(168, 192, 216), true)
	dxlib.DrawBox(x, y, x+120, y+105, dxlib.GetColor(16, 80, 104), true)
	switch state {
	case stateUnsetPartsSelect, stateBoardPartsSelect, stateDeployment:
		c := itemList.GetPointer() + itemList.GetScroll()
		if c < len(unsetParts) {
			parts := ncparts.Get(unsetParts[c].rawData.ID)
			for i, s := range common.SplitMsg(parts.Description, 7) {
				draw.String(x+5, y+5+i*20, 0xFFFFFF, s)
			}
		}
	case stateRun:
		str := "ＲＵＮ・"
		switch count / 20 % 3 {
		case 1:
			str += "・"
		case 2:
			str += "・・"
		}
		draw.String(x+5, y+5, 0xFFFFFF, str)

		// TODO: RUNNING Line
	case stateRunEnd:
		if checkBugs() {
			draw.String(x+5, y+5, 0xFFFFFF, "OK!")
			draw.String(x+5, y+25, 0xFFFFFF, "異常なし")
		} else {
			draw.String(x+5, y+5, 0xFFFFFF, "異常発生")
			draw.String(x+5, y+25, 0xFFFFFF, "プログラムを見直してください")
		}
	}
}

func Process() bool {
	switch state {
	case stateOpening:
		if count == 0 {
			fade.In(30)
		}

		if count > 30 {
			stateChange(stateUnsetPartsSelect)
		}
	case stateUnsetPartsSelect:
		if inputs.CheckKey(inputs.KeyCancel) == 1 {
			sound.On(resources.SECancel)
			return true
		}

		selected = itemList.Process()
		if selected != -1 {
			sound.On(resources.SEMenuEnter)
			if itemList.GetList()[selected] == runName {
				stateChange(stateRun)
			} else {
				stateChange(stateDeployment)
			}
			return false
		}

		if inputs.CheckKey(inputs.KeyLeft) == 1 {
			stateChange(stateBoardPartsSelect)
			return false
		}
	case stateBoardPartsSelect:
		// TODO
	case stateDeployment:
		if inputs.CheckKey(inputs.KeyCancel) == 1 {
			sound.On(resources.SECancel)
			selected = -1
			stateChange(beforeState)
			return false
		}

		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			if checkSet() {
				sound.On(resources.SEMenuEnter)
				newInfo := unsetParts[selected]
				newInfo.rawData.IsSet = true
				newInfo.rawData.X = setPointerPos.X
				newInfo.rawData.Y = setPointerPos.Y

				updateParts(newInfo)
				initParts()
				selected = -1
				stateChange(stateUnsetPartsSelect)
			} else {
				sound.On(resources.SEDenied)
			}
			return false
		}

		if inputs.CheckKey(inputs.KeyUp) == 1 && setPointerPos.Y > 0 {
			sound.On(resources.SECursorMove)
			setPointerPos.Y--
			return false
		}
		if inputs.CheckKey(inputs.KeyDown) == 1 && setPointerPos.Y < boardSize-1 {
			sound.On(resources.SECursorMove)
			setPointerPos.Y++
			return false
		}
		if inputs.CheckKey(inputs.KeyLeft) == 1 && setPointerPos.X > 0 {
			sound.On(resources.SECursorMove)
			setPointerPos.X--
			return false
		}
		if inputs.CheckKey(inputs.KeyRight) == 1 && setPointerPos.X < boardSize-1 {
			sound.On(resources.SECursorMove)
			setPointerPos.X++
			return false
		}
	case stateRun:
		if count >= 30 {
			stateChange(stateRunEnd)
			return false
		}
	case stateRunEnd:
		if count == 0 {
			sound.On(resources.SERunOK)
			parts := []player.NaviCustomParts{}
			for _, p := range allParts {
				parts = append(parts, p.rawData)
			}

			playerInfo.SetNaviCustomParts(parts)
		}

		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			if checkBugs() {
				return true
			} else {
				selected = -1
				stateChange(stateUnsetPartsSelect)
				return false
			}
		}
	}

	count++
	return false
}

func stateChange(next int) {
	logger.Info("Change navu cutom state from %d to %d", state, next)
	beforeState = state
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
	case ncparts.ColorWhite:
		return 0
	case ncparts.ColorYellow:
		return 1
	case ncparts.ColorPink:
		return 2
	}

	common.SetError(fmt.Sprintf("カラーコード %d に対するブロックは存在しません", color))
	return 0
}

func drawBoardParts(basePos common.Point, parts ncparts.NaviCustomParts) {
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
			dxlib.DrawBox(b.X*40+baseX+4, b.Y*40+baseY+4, (b.X+1)*40+baseX, (b.Y+1)*40+baseY, ncparts.GetColorCode(parts.Color), true)
		}
	}
}

func checkSet() bool {
	parts := ncparts.Get(unsetParts[selected].rawData.ID)

	for _, b := range parts.Blocks {
		x := setPointerPos.X + b.X
		y := setPointerPos.Y + b.Y
		// セットするパートがボード外にはみ出ていないか
		if x < 0 || x >= boardSize || y < 0 || y >= boardSize {
			return false
		}

		// パーツ同士が重なっていないか
		for _, s := range setParts {
			setParts := ncparts.Get(s.rawData.ID)
			for _, sb := range setParts.Blocks {
				sx := s.rawData.X + sb.X
				sy := s.rawData.Y + sb.Y
				if x == sx && y == sy {
					return false
				}
			}
		}
	}

	return true // セットできる
}

func initParts() {
	unsetParts = []partsInfo{}
	setParts = []partsInfo{}
	names := []string{}
	for _, p := range allParts {
		if !p.rawData.IsSet {
			unsetParts = append(unsetParts, p)
			parts := ncparts.Get(p.rawData.ID)
			names = append(names, parts.Name)
		} else {
			setParts = append(setParts, p)
		}
	}
	names = append(names, runName)
	itemList.SetList(names, maxListNum)
}

func updateParts(parts partsInfo) {
	for i := range allParts {
		if allParts[i].objID == parts.objID {
			allParts[i] = parts
			return
		}
	}
}

func checkBugs() bool {
	// ルール
	// 　- Plusパーツがライン上にある
	// 　- プログラムパーツがライン上にない
	// 　- (未実装)同じ色のプログラムやプラスパーツは隣同士に置いてはならない
	// 　- (未実装)組み込めるパーツは最大4色まで

	for _, p := range setParts {
		parts := ncparts.Get(p.rawData.ID)
		if parts.IsPlusParts {
			for _, b := range parts.Blocks {
				if p.rawData.Y+b.Y == lineY {
					return false
				}
			}
		} else {
			ok := false
			for _, b := range parts.Blocks {
				if p.rawData.Y+b.Y == lineY {
					ok = true
					break
				}
			}
			if !ok {
				return false
			}
		}
	}

	return true // ok
}
