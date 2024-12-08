package fade

import (
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	workNone int = iota
	workFadeIn
	workFadeOut
)

var (
	workType int
	bright   int
	speed    int
)

func In(count int) {
	workType = workFadeIn
	bright = 0
	speed = 255 / count
}

func Draw() {
	switch workType {
	case workFadeIn, workFadeOut:
		dxlib.SetDrawBright(bright, bright, bright)
	}
}

func Update() {
	switch workType {
	case workFadeIn:
		bright += speed
		if bright > 255 {
			bright = 255
			workType = workNone
		}
	}
}
