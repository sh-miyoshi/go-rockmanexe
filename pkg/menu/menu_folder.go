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

var (
	imgChipFrame  int32
	imgPointer    int32
	playerInfo    *player.Player
	folderPointer = 0
	folderScroll  = 0
)

func folderInit(plyr *player.Player) error {
	playerInfo = plyr
	folderPointer = 0
	folderScroll = 0

	fname := common.ImagePath + "menu/chip_frame.png"
	imgChipFrame = dxlib.LoadGraph(fname)
	if imgChipFrame == -1 {
		return fmt.Errorf("failed to load menu chip frame image %s", fname)
	}

	fname = common.ImagePath + "menu/pointer.png"
	imgPointer = dxlib.LoadGraph(fname)
	if imgPointer == -1 {
		return fmt.Errorf("failed to load menu pointer image %s", fname)
	}

	return nil
}

func folderEnd() {
	dxlib.DeleteGraph(imgChipFrame)
}

func folderProcess() {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		// TODO chip exchange
		sound.On(sound.SEDenied)
	} else if inputs.CheckKey(inputs.KeyCancel) == 1 {
		stateChange(stateTop)
	} else {
		if inputs.CheckKey(inputs.KeyUp)%10 == 1 {
			if folderPointer > 0 {
				sound.On(sound.SESelect)
				folderPointer--
			} else if folderScroll > 0 {
				sound.On(sound.SESelect)
				folderScroll--
			}
		} else if inputs.CheckKey(inputs.KeyDown)%10 == 1 {
			if folderPointer < folderShowNum-1 {
				sound.On(sound.SESelect)
				folderPointer++
			} else if folderScroll < player.FolderSize-folderShowNum {
				sound.On(sound.SESelect)
				folderScroll++
			}
		}
	}
}

func folderDraw() {
	draw.String(20, 5, 0xffffff, "Chip Folder")

	dxlib.DrawBox(25, 55, 460, 300, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawBox(45, 35, 325, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawTriangle(25, 55, 45, 35, 45, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawTriangle(325, 55, 325, 35, 345, 55, dxlib.GetColor(168, 192, 216), dxlib.TRUE)
	dxlib.DrawBox(215, 70, 440, 290, dxlib.GetColor(16, 80, 104), dxlib.TRUE)

	dxlib.DrawGraph(35, 40, imgChipFrame, dxlib.TRUE)
	draw.String(220, 40, 0xffffff, "フォルダ")

	// Show chip list
	for i := 0; i < folderShowNum; i++ {
		n := i + folderScroll
		if n >= player.FolderSize {
			break
		}

		c := playerInfo.ChipFolder[n]
		info := chip.Get(c.ID)
		y := 75 + int32(i)*30
		dxlib.DrawGraph(220, y, chip.GetIcon(c.ID, true), dxlib.TRUE)
		draw.String(250, y+5, 0xffffff, info.Name)
		dxlib.DrawGraph(380, y, chip.GetTypeImage(info.Type), dxlib.TRUE)
		// TODO font
		draw.String(412, y+5, 0xffffff, strings.ToUpper(c.Code))
	}

	// Show pointer
	dxlib.DrawGraph(200, 78+int32(folderPointer)*30, imgPointer, dxlib.TRUE)

	// Show pointered chip detail
	c := playerInfo.ChipFolder[folderPointer+folderScroll]
	info := chip.Get(c.ID)
	dxlib.DrawGraph(58, 65, info.Image, dxlib.TRUE)
	draw.ChipCode(50, 165, c.Code, 100)
	dxlib.DrawGraph(85, 165, chip.GetTypeImage(info.Type), dxlib.TRUE)
	if info.Power > 0 {
		draw.Number(115, 165, int32(info.Power), draw.NumberOption{
			RightAligned: true,
			Length:       4,
		})
	}
}
