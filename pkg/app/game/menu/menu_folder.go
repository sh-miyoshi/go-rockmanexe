package menu

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip"
	chipimage "github.com/sh-miyoshi/go-rockmanexe/pkg/app/chip/image"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/player"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/window"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
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
	imgChipFrame     int
	imgArrow         int
	imgPointer       int
	imgScrollPointer int

	pointer       [folderWindowTypeMax]int
	scroll        [folderWindowTypeMax]int
	listNum       [folderWindowTypeMax]int
	selected      int
	currentWindow int
	count         int
	win           window.MessageWindow

	msg string

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
	fname := config.ImagePath + "menu/chip_frame.png"
	res.imgChipFrame = dxlib.LoadGraph(fname)
	if res.imgChipFrame == -1 {
		return nil, fmt.Errorf("failed to load menu chip frame image %s", fname)
	}

	fname = config.ImagePath + "menu/pointer.png"
	res.imgPointer = dxlib.LoadGraph(fname)
	if res.imgPointer == -1 {
		return nil, fmt.Errorf("failed to load menu pointer image %s", fname)
	}

	fname = config.ImagePath + "menu/arrow.png"
	res.imgArrow = dxlib.LoadGraph(fname)
	if res.imgArrow == -1 {
		return nil, fmt.Errorf("failed to load menu arrow image %s", fname)
	}

	fname = config.ImagePath + "menu/scroll_point.png"
	res.imgScrollPointer = dxlib.LoadGraph(fname)
	if res.imgScrollPointer == -1 {
		return nil, fmt.Errorf("failed to load menu scroll point image %s", fname)
	}

	if err := res.win.Init(); err != nil {
		return nil, err
	}

	return res, nil
}

func (f *menuFolder) End() {
	dxlib.DeleteGraph(f.imgChipFrame)
	dxlib.DeleteGraph(f.imgPointer)
	dxlib.DeleteGraph(f.imgArrow)
	dxlib.DeleteGraph(f.imgScrollPointer)
}

func (f *menuFolder) Process() {
	if f.msg != "" {
		if f.win.Process() {
			sound.On(resources.SECancel)
			f.msg = ""
		}
		return
	}

	f.count++

	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		c := f.currentWindow
		sel := c*player.FolderSize + f.scroll[c] + f.pointer[c]

		if f.currentWindow == folderWindowTypeBackPack && f.listNum[f.currentWindow] == 0 {
			sound.On(resources.SEDenied)
			return
		}

		if f.selected == -1 {
			// First select
			f.selected = sel
			sound.On(resources.SESelect)
		} else {
			if err := f.exchange(f.selected, sel); err != nil {
				sound.On(resources.SEDenied)
				logger.Info("Failed to exchange chip: %v", err)
				f.msg = err.Error()
				f.win.SetMessage(f.msg, window.FaceTypeNone)
				return
			}

			sound.On(resources.SESelect)
			f.selected = -1
		}
	} else if inputs.CheckKey(inputs.KeyCancel) == 1 {
		if f.selected == -1 {
			stateChange(stateTop)
		} else {
			f.selected = -1
		}
		sound.On(resources.SECancel)
	} else {
		if inputs.CheckKey(inputs.KeyUp)%10 == 1 {
			if f.pointer[f.currentWindow] > 0 {
				sound.On(resources.SECursorMove)
				f.pointer[f.currentWindow]--
			} else if f.scroll[f.currentWindow] > 0 {
				sound.On(resources.SECursorMove)
				f.scroll[f.currentWindow]--
			}
		} else if inputs.CheckKey(inputs.KeyDown)%10 == 1 {
			n := folderShowNum - 1
			if f.listNum[f.currentWindow] < folderShowNum {
				n = f.listNum[f.currentWindow] - 1
			}

			if f.pointer[f.currentWindow] < n {
				sound.On(resources.SECursorMove)
				f.pointer[f.currentWindow]++
			} else if f.scroll[f.currentWindow] < f.listNum[f.currentWindow]-folderShowNum {
				sound.On(resources.SECursorMove)
				f.scroll[f.currentWindow]++
			}
		} else if inputs.CheckKey(inputs.KeyLeft) == 1 {
			if f.currentWindow == folderWindowTypeBackPack {
				f.currentWindow = folderWindowTypeFolder
				sound.On(resources.SEWindowChange)
			}
		} else if inputs.CheckKey(inputs.KeyRight) == 1 {
			if f.currentWindow == folderWindowTypeFolder {
				f.currentWindow = folderWindowTypeBackPack
				sound.On(resources.SEWindowChange)
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
		x := 220
		if f.currentWindow == folderWindowTypeBackPack {
			x = 50
		}

		y := 75 + i*30
		dxlib.DrawGraph(x, y, chipimage.GetIcon(c.ID, true), true)
		draw.String(x+30, y+5, 0xffffff, info.Name)
		dxlib.DrawGraph(x+160, y, chipimage.GetType(info.Type), true)
		// TODO font
		draw.String(x+192, y+5, 0xffffff, strings.ToUpper(c.Code))
	}

	// Show pointer
	if f.listNum[f.currentWindow] > 0 {
		tx := 200
		if f.currentWindow == folderWindowTypeBackPack {
			tx = 30
		}

		if f.selected != -1 && (f.count/2)%2 == 0 {
			win := f.selected / player.FolderSize
			if win == f.currentWindow {
				sel := f.selected
				if sel >= player.FolderSize {
					sel -= player.FolderSize
				}

				areaBegin := f.scroll[f.currentWindow]
				areaEnd := f.scroll[f.currentWindow] + folderShowNum
				if areaBegin <= sel && sel <= areaEnd {
					p := sel - areaBegin
					dxlib.DrawGraph(tx+2, 80+p*30, f.imgPointer, true)
				}
			}
		}

		dxlib.DrawGraph(tx, 78+f.pointer[f.currentWindow]*30, f.imgPointer, true)
	}

	// Show scrol bar
	var length, start int
	switch f.currentWindow {
	case folderWindowTypeFolder:
		start = 80
		length = 205
	case folderWindowTypeBackPack:
		start = 60
		length = 225
	}
	dxlib.DrawBox(450, start, 453, start+length, dxlib.GetColor(16, 80, 104), true)
	n := f.scroll[f.currentWindow] + f.pointer[f.currentWindow]
	dxlib.DrawGraph(442, start+(length-23)*n/f.listNum[f.currentWindow]-1, f.imgScrollPointer, true)

	// Show pointered chip detail
	f.drawChipDetail(f.pointer[f.currentWindow] + f.scroll[f.currentWindow])

	// Show message
	if f.msg != "" {
		f.win.Draw()
	}
}

func (f *menuFolder) drawBackGround() {
	dxlib.DrawBox(25, 55, 460, 300, dxlib.GetColor(168, 192, 216), true)

	switch f.currentWindow {
	case folderWindowTypeFolder:
		dxlib.DrawBox(45, 35, 325, 55, dxlib.GetColor(168, 192, 216), true)
		dxlib.DrawTriangle(25, 55, 45, 35, 45, 55, dxlib.GetColor(168, 192, 216), true)
		dxlib.DrawTriangle(325, 55, 325, 35, 345, 55, dxlib.GetColor(168, 192, 216), true)
		dxlib.DrawBox(215, 70, 440, 290, dxlib.GetColor(16, 80, 104), true)

		dxlib.DrawGraph(35, 40, f.imgChipFrame, true)
		dxlib.DrawGraph(443, 57, f.imgArrow, true)
		draw.String(220, 40, 0xffffff, "フォルダ")
	case folderWindowTypeBackPack:
		dxlib.DrawBox(440, 35, 160, 55, dxlib.GetColor(168, 192, 216), true)
		dxlib.DrawTriangle(460, 55, 440, 35, 440, 55, dxlib.GetColor(168, 192, 216), true)
		dxlib.DrawTriangle(160, 55, 160, 35, 140, 55, dxlib.GetColor(168, 192, 216), true)
		dxlib.DrawBox(270, 70, 45, 290, dxlib.GetColor(16, 80, 104), true)

		dxlib.DrawGraph(285, 40, f.imgChipFrame, true)
		dxlib.DrawTurnGraph(28, 57, f.imgArrow, true)
		draw.String(200, 40, 0xffffff, "リュック")
	}
}

func (f *menuFolder) drawChipDetail(index int) {
	if index >= f.listNum[f.currentWindow] {
		return
	}

	var c player.ChipInfo
	var ofsx int
	switch f.currentWindow {
	case folderWindowTypeFolder:
		c = f.playerInfo.ChipFolder[index]
		ofsx = 35
	case folderWindowTypeBackPack:
		c = f.playerInfo.BackPack[index]
		ofsx = 285
	}

	info := chip.Get(c.ID)
	dxlib.DrawGraph(23+ofsx, 65, chipimage.GetDetail(c.ID), true)
	draw.ChipCode(15+ofsx, 165, c.Code, 100)
	dxlib.DrawGraph(50+ofsx, 165, chipimage.GetType(info.Type), true)
	if info.Power > 0 {
		draw.Number(80+ofsx, 165, int(info.Power), draw.NumberOption{
			RightAligned: true,
			Length:       4,
		})
	}

	f.drawDescription(info.Description)
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
		id := f.playerInfo.BackPack[backPackSel].ID

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

func (f *menuFolder) drawDescription(desc string) {
	splited := []string{}
	for len(desc) > 0 {
		res := ""
		for i := 0; i < 9 && len(desc) > 0; i++ {
			r, size := utf8.DecodeRuneInString(desc)
			res += string(r)
			desc = desc[size:]
		}
		splited = append(splited, res)
	}

	switch f.currentWindow {
	case folderWindowTypeFolder:
		for i, d := range splited {
			draw.String(50, 205+i*25, 0x000000, d)
		}
	case folderWindowTypeBackPack:
		for i, d := range splited {
			draw.String(300, 205+i*25, 0x000000, d)
		}
	}
}
