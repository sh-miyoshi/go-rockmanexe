package inputs

import "github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"

// KeyType ...
type KeyType int

const (
	KeyEnter KeyType = iota
	KeyCancel
	KeyLeft
	KeyRight
	KeyUp
	KeyDown
	KeyLButton
	KeyRButton
	KeyDebug

	keyMax
)

var (
	keyState [256]int
	keyBind  [keyMax]int
)

// InitByDefault set key binding by default value
func InitByDefault() {
	keyBind[KeyEnter] = dxlib.KEY_INPUT_Z
	keyBind[KeyCancel] = dxlib.KEY_INPUT_X
	keyBind[KeyLeft] = dxlib.KEY_INPUT_LEFT
	keyBind[KeyRight] = dxlib.KEY_INPUT_RIGHT
	keyBind[KeyUp] = dxlib.KEY_INPUT_UP
	keyBind[KeyDown] = dxlib.KEY_INPUT_DOWN
	keyBind[KeyLButton] = dxlib.KEY_INPUT_A
	keyBind[KeyRButton] = dxlib.KEY_INPUT_S
	keyBind[KeyDebug] = dxlib.KEY_INPUT_D
}

// TODO: InitBySetting(settingFile string)

// KeyStateUpdate ...
func KeyStateUpdate() {
	tmp := make([]byte, 256)
	dxlib.GetHitKeyStateAll(tmp)
	for i := 0; i < 256; i++ {
		if tmp[i] == 1 {
			keyState[i]++
		} else {
			keyState[i] = 0
		}
	}
}

// CheckKey ...
func CheckKey(key KeyType) int {
	return keyState[keyBind[key]]
}
