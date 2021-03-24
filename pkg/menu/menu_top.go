package menu

import (
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

const (
	topSelectChipFolder int = iota
	topSelectGoBattle
	topSelectRecord

	topSelectMax
)

var (
	pointer = 0
)

func topInit() {
}

func topEnd() {
}

func topProcess() {
	if inputs.CheckKey(inputs.KeyEnter) == 1 {
		// TODO stateChange()
	} else {
		if inputs.CheckKey(inputs.KeyUp) == 1 {
			if pointer > 0 {
				pointer--
			}
		} else if inputs.CheckKey(inputs.KeyDown) == 1 {
			if pointer < topSelectMax-1 {
				pointer++
			}
		}
	}
}

func topDraw() {
	msgs := [topSelectMax]string{
		"チップフォルダー",
		"バトル",
		"戦績",
	}

	x := int32(40)
	for i, msg := range msgs {
		draw.String(x+20, 50+int32(i)*25, 0xffffff, msg)
	}

	const s = 2
	y := int32(50 + pointer*25)
	dxlib.DrawTriangle(x, y+s, x+18-s*2, y+10, x, y+20-s, 0xffffff, dxlib.TRUE)

	// TODO Show description
}
