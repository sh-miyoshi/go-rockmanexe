package inputs

import "github.com/sh-miyoshi/dxlib"

const padButtonNum = 28

type pad struct {
	keyBind  [keyMax]int
	padState [padButtonNum]int
}

func (p *pad) Init() error {
	p.keyBind[KeyEnter] = 6
	p.keyBind[KeyCancel] = 5
	p.keyBind[KeyLeft] = 2
	p.keyBind[KeyRight] = 3
	p.keyBind[KeyUp] = 4
	p.keyBind[KeyDown] = 1
	p.keyBind[KeyLButton] = 9
	p.keyBind[KeyRButton] = 10
	p.keyBind[KeyDebug] = 12

	return nil
}

func (p *pad) KeyStateUpdate() {
	state := dxlib.GetJoypadInputState(dxlib.DX_INPUT_PAD1 | dxlib.DX_INPUT_KEY)
	for i := 0; i < padButtonNum; i++ {
		if state&(1<<i) != 0 {
			p.padState[i]++
		} else {
			p.padState[i] = 0
		}
	}
}

func (p *pad) CheckKey(key KeyType) int {
	return p.padState[p.keyBind[key]-1]
}
