//go:build mac
// +build mac

package dxlib

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

	panic("not implemented yet")
}

func LoadDivGraph(fname string, allNum, xnum, ynum, xsize, ysize int, handleBuf []int) int {
	if disabled {
		return 0
	}

	panic("not implemented yet")
}

func DeleteGraph(grHandle int) int {
	if disabled {
		return 0
	}

	panic("not implemented yet")
}

func DrawGraph(x, y int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawExtendGraph(x1, y1, x2, y2 int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func CreateFontToHandle(opt ...CreateFontToHandleOption) int {
	if disabled {
		return 0
	}

	panic("not implemented yet")
}

func DrawFormatString(x, y int, color uint, format string, a ...interface{}) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawFormatStringToHandle(x, y int, color uint, fontHandle int, format string, a ...interface{}) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawStringToHandle(x, y int, color uint, fontHandle int, message string) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawExtendFormatStringToHandle(x, y int, exRateX, exRateY float64, color uint, fontHandle int, format string, a ...interface{}) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func SetDrawBlendMode(blendMode int, pal int) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawRotaGraph(x, y int, extRate, angle float64, grHandle int, transFlag bool, opt ...DrawRotaGraphOption) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawBox(x1, y1, x2, y2 int, color uint, fillFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func GetColor(red, green, blue int) uint {
	if disabled {
		return 0
	}

	panic("not implemented yet")
}

func SetDrawBright(redBright, greenBright, blueBright int) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func LoadSoundMem(fname string) int {
	if disabled {
		return 0
	}

	panic("not implemented yet")
}

func PlaySoundMem(soundHandle int, playType int, topPositionFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func CheckSoundMem(soundHandle int) int {
	if disabled {
		return 0
	}

	panic("not implemented yet")
}

func StopSoundMem(soundHandle int) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawTriangle(x1, y1, x2, y2, x3, y3 int, color uint, fillFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func ChangeVolumeSoundMem(volumePan int, soundHandle int) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawTurnGraph(x, y int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func GetDrawStringWidth(str string, strLen int) int {
	if disabled {
		return 0
	}

	panic("not implemented yet")
}

func GetHitKeyStateAll(keyStateBuf []byte) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func GetGraphSize(grHandle int, sizeX, sizeY *int) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawRectGraph(destX, destY, srcX, srcY int, width, height int, grHandle int, transFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawCircle(x, y int, r int, color uint, fillFlag bool) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func DrawLine(x1, y1, x2, y2 int, color uint) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func SetWindowSize(width int32, height int32) {
	if disabled {
		return
	}

	panic("not implemented yet")
}

func GetJoypadInputState(inputType int32) int32 {
	panic("not implemented yet")
}
