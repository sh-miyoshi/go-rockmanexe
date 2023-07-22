//go:build !mac
// +build !mac

package dxlib

import (
	"github.com/sh-miyoshi/dxlib"
)

var (
	disabled = false
)

func Disable() {
	disabled = true
}

func Int32Ptr(a int32) *int32 {
	return &a
}

func StringPtr(a string) *string {
	return &a
}

func LoadGraph(fname string) int {
	if disabled {
		return 0
	}

	return int(dxlib.LoadGraph(fname))
}

func LoadDivGraph(fname string, allNum, xnum, ynum, xsize, ysize int, handleBuf []int) int {
	if disabled {
		return 0
	}

	tmpBuf := make([]int32, len(handleBuf))
	res := dxlib.LoadDivGraph(fname, int32(allNum), int32(xnum), int32(ynum), int32(xsize), int32(ysize), tmpBuf)
	for i := 0; i < len(handleBuf); i++ {
		handleBuf[i] = int(tmpBuf[i])
	}

	return int(res)
}

func DeleteGraph(grHandle int) int {
	if disabled {
		return 0
	}

	return int(dxlib.DeleteGraph(int32(grHandle)))
}

func DrawGraph(x, y int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	dxlib.DrawGraph(int32(x), int32(y), int32(grHandle), makeFlag(transFlag))
}

func DrawExtendGraph(x1, y1, x2, y2 int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	dxlib.DrawExtendGraph(int32(x1), int32(y1), int32(x2), int32(y2), int32(grHandle), makeFlag(transFlag))
}

func CreateFontToHandle(opt ...CreateFontToHandleOption) int {
	if disabled {
		return 0
	}

	if opt != nil {
		dxopt := dxlib.CreateFontToHandleOption{
			FontName: opt[0].FontName,
			Size:     opt[0].Size,
			Thick:    opt[0].Thick,
			FontType: opt[0].FontType,
			CharSet:  opt[0].CharSet,
			EdgeSize: opt[0].EdgeSize,
			Italic:   opt[0].Italic,
			Handle:   opt[0].Handle,
		}
		return int(dxlib.CreateFontToHandle(dxopt))
	}
	return int(dxlib.CreateFontToHandle())
}

func DrawFormatString(x, y int, color uint, format string, a ...interface{}) {
	if disabled {
		return
	}

	dxlib.DrawFormatString(int32(x), int32(y), uint32(color), format, a...)
}

func DrawFormatStringToHandle(x, y int, color uint, fontHandle int, format string, a ...interface{}) {
	if disabled {
		return
	}

	dxlib.DrawFormatStringToHandle(int32(x), int32(y), uint32(color), int32(fontHandle), format, a...)
}

func DrawStringToHandle(x, y int, color uint, fontHandle int, message string) {
	if disabled {
		return
	}

	dxlib.DrawStringToHandle(int32(x), int32(y), message, uint32(color), int32(fontHandle))
}

func DrawExtendFormatStringToHandle(x, y int, exRateX, exRateY float64, color uint, fontHandle int, format string, a ...interface{}) {
	if disabled {
		return
	}

	dxlib.DrawExtendFormatStringToHandle(int32(x), int32(y), exRateX, exRateY, uint32(color), int32(fontHandle), format, a...)
}

func SetDrawBlendMode(blendMode int, pal int) {
	if disabled {
		return
	}

	dxlib.SetDrawBlendMode(int32(blendMode), int32(pal))
}

func DrawRotaGraph(x, y int, extRate, angle float64, grHandle int, transFlag bool, opt ...DrawRotaGraphOption) {
	if disabled {
		return
	}

	if opt != nil {
		dxopt := dxlib.DrawRotaGraphOption{
			ReverseXFlag: opt[0].ReverseXFlag,
			ReverseYFlag: opt[0].ReverseYFlag,
		}
		dxlib.DrawRotaGraph(int32(x), int32(y), extRate, angle, int32(grHandle), makeFlag(transFlag), dxopt)
	} else {
		dxlib.DrawRotaGraph(int32(x), int32(y), extRate, angle, int32(grHandle), makeFlag(transFlag))
	}
}

func DrawBox(x1, y1, x2, y2 int, color uint, fillFlag bool) {
	if disabled {
		return
	}

	dxlib.DrawBox(int32(x1), int32(y1), int32(x2), int32(y2), uint32(color), makeFlag(fillFlag))
}

func GetColor(red, green, blue int) uint {
	if disabled {
		return 0
	}

	return uint(dxlib.GetColor(int32(red), int32(green), int32(blue)))
}

func SetDrawBright(redBright, greenBright, blueBright int) {
	if disabled {
		return
	}

	dxlib.SetDrawBright(int32(redBright), int32(greenBright), int32(blueBright))
}

func LoadSoundMem(fname string) int {
	if disabled {
		return 0
	}

	return int(dxlib.LoadSoundMem(fname))
}

func PlaySoundMem(soundHandle int, playType int, topPositionFlag bool) {
	if disabled {
		return
	}

	dxlib.PlaySoundMem(int32(soundHandle), int32(playType), makeFlag(topPositionFlag))
}

func CheckSoundMem(soundHandle int) int {
	if disabled {
		return 0
	}

	return int(dxlib.CheckSoundMem(int32(soundHandle)))
}

func StopSoundMem(soundHandle int) {
	if disabled {
		return
	}

	dxlib.StopSoundMem(int32(soundHandle))
}

func DrawTriangle(x1, y1, x2, y2, x3, y3 int, color uint, fillFlag bool) {
	if disabled {
		return
	}

	dxlib.DrawTriangle(int32(x1), int32(y1), int32(x2), int32(y2), int32(x3), int32(y3), uint32(color), makeFlag(fillFlag))
}

func makeFlag(flag bool) int32 {
	if flag {
		return dxlib.TRUE
	} else {
		return dxlib.FALSE
	}
}

func ChangeVolumeSoundMem(volumePan int, soundHandle int) {
	if disabled {
		return
	}

	dxlib.ChangeVolumeSoundMem(int32(volumePan), int32(soundHandle))
}

func DrawTurnGraph(x, y int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	dxlib.DrawTurnGraph(int32(x), int32(y), int32(grHandle), makeFlag(transFlag))
}

func GetDrawStringWidth(str string, strLen int) int {
	if disabled {
		return 0
	}

	return int(dxlib.GetDrawStringWidth(str, int32(strLen)))
}

func GetHitKeyStateAll(keyStateBuf []byte) {
	if disabled {
		return
	}

	dxlib.GetHitKeyStateAll(keyStateBuf)
}

func GetGraphSize(grHandle int, sizeX, sizeY *int) {
	if disabled {
		return
	}

	var tx, ty int32
	dxlib.GetGraphSize(int32(grHandle), &tx, &ty)
	*sizeX = int(tx)
	*sizeY = int(ty)
}

func DrawRectGraph(destX, destY, srcX, srcY int, width, height int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	dxlib.DrawRectGraph(int32(destX), int32(destY), int32(srcX), int32(srcY), int32(width), int32(height), int32(grHandle), makeFlag(transFlag))
}

func DrawCircle(x, y int, r int, color uint, fillFlag bool) {
	if disabled {
		return
	}

	dxlib.DrawCircle(int32(x), int32(y), int32(r), uint32(color), makeFlag(fillFlag))
}

func DrawLine(x1, y1, x2, y2 int, color uint) {
	if disabled {
		return
	}

	dxlib.DrawLine(int32(x1), int32(y1), int32(x2), int32(y2), uint32(color))
}

func SetWindowSize(width int32, height int32) {
	if disabled {
		return
	}

	dxlib.SetWindowSize(width, height)
}
