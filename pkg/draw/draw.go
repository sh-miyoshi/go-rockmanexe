package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

var (
	fontHandle int32 = -1
	imgCode    []int32
)

func Init() error {
	// Set font
	fontHandle = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
		FontName: dxlib.StringPtr("k8x12"),
		Size:     dxlib.Int32Ptr(22),
		Thick:    dxlib.Int32Ptr(7),
	})
	if fontHandle == -1 {
		return fmt.Errorf("Failed to create font")
	}

	// Load chip code
	imgCode = make([]int32, 27)
	fname := common.ImagePath + "chip_code.png"
	if res := dxlib.LoadDivGraph(fname, 27, 9, 3, 20, 26, imgCode); res == -1 {
		return fmt.Errorf("Failed to load chip code image %s", fname)
	}

	// TODO: load number data

	return nil
}

func String(x int32, y int32, color uint32, format string, a ...interface{}) {
	dxlib.DrawFormatStringToHandle(x, y, color, fontHandle, format, a...)
}

func ChipCode(x int32, y int32, code string) {
	index := -1
	if len(code) != 1 {
		panic(fmt.Sprintf("Invalid chip code %s is specified.", code))
	}

	rc := []rune(code)
	if rc[0] >= 'a' && rc[0] <= 'z' {
		index = int(rc[0] - 'a')
	} else if rc[0] >= 'A' && rc[0] <= 'Z' {
		index = int(rc[0] - 'A')
	} else if rc[0] == '*' {
		index = 26
	} else {
		panic(fmt.Sprintf("Invalid chip code %s is specified.", code))
	}

	dxlib.DrawGraph(x, y, imgCode[index], dxlib.FALSE)
}
