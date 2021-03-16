package draw

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
)

var (
	fontHandle int32 = -1
)

func Init() error {
	fontHandle = dxlib.CreateFontToHandle(dxlib.CreateFontToHandleOption{
		FontName: dxlib.StringPtr("k8x12"),
		Size:     dxlib.Int32Ptr(22),
		Thick:    dxlib.Int32Ptr(7),
	})
	if fontHandle == -1 {
		return fmt.Errorf("Failed to create font")
	}

	// TODO: load number data

	return nil
}

func String(x int32, y int32, color uint32, format string, a ...interface{}) {
	dxlib.DrawFormatStringToHandle(x, y, color, fontHandle, format, a...)
}
