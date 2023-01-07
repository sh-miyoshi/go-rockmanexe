package inputs

import "github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"

type keyboard struct {
	keyState [256]int
	keyBind  [keyMax]int
}

func (k *keyboard) Init() error {
	k.keyBind[KeyEnter] = dxlib.KEY_INPUT_Z
	k.keyBind[KeyCancel] = dxlib.KEY_INPUT_X
	k.keyBind[KeyLeft] = dxlib.KEY_INPUT_LEFT
	k.keyBind[KeyRight] = dxlib.KEY_INPUT_RIGHT
	k.keyBind[KeyUp] = dxlib.KEY_INPUT_UP
	k.keyBind[KeyDown] = dxlib.KEY_INPUT_DOWN
	k.keyBind[KeyLButton] = dxlib.KEY_INPUT_A
	k.keyBind[KeyRButton] = dxlib.KEY_INPUT_S
	k.keyBind[KeyDebug] = dxlib.KEY_INPUT_D

	return nil
}

func (k *keyboard) KeyStateUpdate() {
	tmp := make([]byte, 256)
	dxlib.GetHitKeyStateAll(tmp)
	for i := 0; i < 256; i++ {
		if tmp[i] == 1 {
			k.keyState[i]++
		} else {
			k.keyState[i] = 0
		}
	}
}

func (k *keyboard) CheckKey(key KeyType) int {
	return k.keyState[k.keyBind[key]]
}
