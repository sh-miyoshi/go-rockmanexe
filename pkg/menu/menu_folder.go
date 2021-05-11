package menu

import (
	"fmt"
	"strings"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/chip"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/sound"
)

const (
	folderShowNum = 7
)

const (
	folderWindowTypeFolder int = iota
	folderWindowTypeBackPack
	folderWindowTypeMax
)

type menuFolder struct {
	imgChipFrame int32
	imgArrow     int32
	imgPointer   int32

	pointer       [folderWindowTypeMax]int
	scroll        [folderWindowTypeMax]int
	listNum       [folderWindowTypeMax]int
	selected      int
	currentWindow int
	count         int

	playerInfo *player.Player
}

func folderNew(plyr *player.Player) (*menuFolder, error) {
	res := &menuFolder{
		playerInfo:    plyr,
		count:         0,
		selected:      -1,
		currentWindow: folderWindowTypeFolder,
	}

	res.listNum[folderWindowTypeFolder] = player.FolderSize
	res.listNum[folderWindowTypeBackPack] = len(plyr.BackPack)

	// Load images
	fname := common.ImagePath + "menu/chip_frame.png"
	res.imgChipFrame = dxlib.LoadGraph(fname)
	if res.imgChipFrame == -1 {
		return nil, fmt.Errorf("failed to load menu chip frame image %s", fname)
	}

	fname = common.ImagePath + "menu/pointer.png"
	res.imgPointer = dxlib.LoadGraph(fname)
	if res.imgPointer == -1 {
		return nil, fmt.Errorf("failed to load menu pointer image %s", fname)
	}

	fname = common.ImagePath + "menu/arrow.png"
	res.imgArrow = dxlib.LoadGraph(fname)
	if res.imgArrow == -1 {
		return nil, fmt.Errorf("failed to load menu arrow image %s", fname)
	}

	return res, nil
}

func (f *menuFolder) End() {
	dxlib.DeleteGraph(f.imgChipFrame)
	dxlib.DeleteGraph(f.imgPointer)
	dxlib.DeleteGraph(f.imgArrow)
}

func (f *menuFolder) Process() {
	f.count++

	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		c := f.currentWindow
		sel := c*player.FolderSize + f.scroll[c]*folderShowNum + f.pointer[c]

		if f.currentWindow == folderWindowTypeBackPack && f.listNum[f.currentWindow] == 0 {
			sound.On(sound.SEDenied)
			return
		}

		if f.selected == -1 {
			// First select
			f.selected = sel
			sound.On(sound.SESelect)
		} else {
			if err := f.exchange(f.selected, sel); err != nil {
				// TODO error handling
				sound.On(sound.SEDenied)
				panic(err)
			}

			sound.On(sound.SESelect)
			f.selected = -1
		}
	} else if inputs.CheckKey(inputs.KeyCancel) == 1 {
		if f.selected == -1 {
			stateChange(stateTop)
		} else {
			f.selected = -1
		}
		sound.On(sound.SECancel)
	} else {
		if inputs.CheckKey(inputs.KeyUp)%10 == 1 {
			if f.pointer[f.currentWindow] > 0 {
				sound.On(sound.SECursorMove)
				f.pointer[f.currentWindow]--
			} else if f.scroll[f.currentWindow] > 0 {
				sound.On(sound.SECursorMove)
				f.scroll[f.currentWindow]--
			}
		} else if inputs.CheckKey(inputs.KeyDown)%10 == 1 {
			n := folderShowNum - 1
			if f.listNum[f.currentWindow] < folderShowNum {
				n = f.listNum[f.currentWindow] - 1
			}

			if f.pointer[f.currentWindow] < n {
				sound.On(sound.SECursorMove)
				f.pointer[f.currentWindow]++
			} else if f.scroll[f.currentWindow] < f.listNum[f.currentWindow]-folderShowNum {
				sound.On(sound.SECursorMove)
				f.scroll[f.currentWindow]++
			}
		} else if inputs.CheckKey(inputs.KeyLeft) == 1 {
			if f.currentWindow == folderWindowTypeBackPack {
				f.currentWindow = folderWindowTypeFolder
				sound.On(sound.SEWindowChange)
			}
		} else if inputs.CheckKey(inputs.KeyRight) == 1 {
			if f.currentWindow == folderWindowTypeFolder {
				f.currentWindow = folderWindowTypeBackPack
				sound.On(sound.SEWindowChange)
			}
		}
	}
}

func (f *menuFolder) Draw() {
	// Show title
	draw.String(20, 5, 0xffffff, "Chip Folder")

	// Show Background
	f.drawBackGround()

	// Show chip list
	for i := 0; i < folderShowNum; i++ {
		n := i + f.scroll[f.currentWindow]
		if n >= f.listNum[f.currentWindow] {
			break
		}

		var c player.ChipInfo
		switch f.currentWindow {
		case folderWindowTypeFolder:
			c = f.playerInfo.ChipFolder[n]
		case folderWindowTypeBackPack:
			c = f.playerInfo.BackPack[n]
		}

		info := chip.Get(c.ID)
		x := int32(220)
		if f.currentWindow == folderWindowTypeBackPack {
			x = 50
		}

		y := 75 + int32(i)*30
		dxlib.DrawGraph(x, y, chip.GetIcon(c.ID, true), dxlib.TRUE)
		draw.String(x+30, y+5, 0xffffff, info.Name)
		dxlib.DrawGraph(x+160, y, chip.GetTypeImage(info.Type), dxlib.TRUE)
		// TODO font
		draw.String(x+192, y+5, 0xffffff, strings.ToUpper(c.Code))
	}

	// Show pointer
	if f.listNum[f.currentWindow] > 0 {
		tx := int32(200)
		if f.currentWindow == folderWindowTypeBackPack {
			tx = 30
		}
		dxlib.DrawGraph(tx, 78+int32(f.pointer[f.currentWindow])*30, f.imgPointer, dxlib.TRUE)
	}

	// Show pointered chip detail
	f.drawChipDetail(f.pointer[f.currentWindow] + f.scroll[f.currentWindow])
}

func (f *menuFolder) drawBackGround() {
	dxlib.DrawBox(25, 55, 460, 300, dxlib.GetColor(168, 192, 216), dxlib.TRUE)

	switch f.currentWindow {
	case folderWindowTypeFolder:
		dxlib.DrawBox(45, 35, 325, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
		dxlib.DrawTriangle(25, 55, 45, 35, 45, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
		dxlib.DrawTriangle(325, 55, 325, 35, 345, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
		dxlib.DrawBox(215, 70, 440, 290, dxlib.GetColor(16, 80, 104), dxlib.TRUE)

		dxlib.DrawGraph(35, 40, f.imgChipFrame, dxlib.TRUE)
		dxlib.DrawGraph(443, 57, f.imgArrow, dxlib.TRUE)
		draw.String(220, 40, 0xffffff, "フォルダ")
	case folderWindowTypeBackPack:
		dxlib.DrawBox(440, 35, 160, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
		dxlib.DrawTriangle(460, 55, 440, 35, 440, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
		dxlib.DrawTriangle(160, 55, 160, 35, 140, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
		dxlib.DrawBox(270, 70, 45, 290, dxlib.GetColor(16, 80, 104), dxlib.TRUE)

		dxlib.DrawGraph(285, 40, f.imgChipFrame, dxlib.TRUE)
		dxlib.DrawTurnGraph(28, 57, f.imgArrow, dxlib.TRUE)
		draw.String(200, 40, 0xffffff, "リュック")
	}
}

func (f *menuFolder) drawChipDetail(index int) {
	if index >= f.listNum[f.currentWindow] {
		return
	}

	var c player.ChipInfo
	var ofsx int32
	switch f.currentWindow {
	case folderWindowTypeFolder:
		c = f.playerInfo.ChipFolder[index]
		ofsx = 35
	case folderWindowTypeBackPack:
		c = f.playerInfo.BackPack[index]
		ofsx = 285
	}

	info := chip.Get(c.ID)
	dxlib.DrawGraph(23+ofsx, 65, info.Image, dxlib.TRUE)
	draw.ChipCode(15+ofsx, 165, c.Code, 100)
	dxlib.DrawGraph(50+ofsx, 165, chip.GetTypeImage(info.Type), dxlib.TRUE)
	if info.Power > 0 {
		draw.Number(80+ofsx, 165, int32(info.Power), draw.NumberOption{
			RightAligned: true,
			Length:       4,
		})
	}
}

func (f *menuFolder) exchange(sel1, sel2 int) error {
	var target1, target2 *player.ChipInfo
	t := 1
	folderSel := 0
	backPackSel := 0

	if sel1 >= player.FolderSize {
		target1 = &f.playerInfo.BackPack[sel1-player.FolderSize]
		backPackSel = sel1 - player.FolderSize
		t *= -1
	} else {
		target1 = &f.playerInfo.ChipFolder[sel1]
		folderSel = sel1
	}

	if sel2 >= player.FolderSize {
		target2 = &f.playerInfo.BackPack[sel2-player.FolderSize]
		backPackSel = sel2 - player.FolderSize
		t *= -1
	} else {
		target2 = &f.playerInfo.ChipFolder[sel2]
		folderSel = sel2
	}

	// Validation
	if t < 0 { // 片方がFolderで、もう片方がBackPackなら
		// Check the number of same name chips in folder
		n := 0
		id := f.playerInfo.ChipFolder[backPackSel].ID

		for i := 0; i < player.FolderSize; i++ {
			if i == folderSel {
				continue
			}
			if f.playerInfo.ChipFolder[i].ID == id {
				n++
			}
		}

		if n >= player.SameChipNumInFolder {
			return fmt.Errorf("同名チップは%d枚までしか入れられません", player.SameChipNumInFolder)
		}
	}

	*target1, *target2 = *target2, *target1
	return nil
}
