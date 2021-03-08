package inputs

import "github.com/sh-miyoshi/dxlib"

var (
	keyState [256]int
)

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
func CheckKey(key int) int {
	return keyState[key]
}
