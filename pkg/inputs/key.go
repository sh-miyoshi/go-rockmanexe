package inputs

import "github.com/sh-miyoshi/dxlib"

// KeyType ...
type KeyType int

const (
	// KeyEnter ...
	KeyEnter KeyType = iota
	// KeyCancel ...
	KeyCancel
	// KeyLeft ...
	KeyLeft
	// KeyRight ...
	KeyRight
	// KeyUp ...
	KeyUp
	// KeyDown ...
	KeyDown

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
