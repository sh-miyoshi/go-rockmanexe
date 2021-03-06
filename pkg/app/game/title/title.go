package title

import (
	"errors"
	"fmt"
	"os"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/inputs"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	stateBegin int = iota
	stateSelect
)

var (
	ErrStartInit     = errors.New("start with initialize")
	ErrStartContinue = errors.New("start with continue")

	count     int
	imgLogo   int32
	imgBack   int32
	state     int
	cursor    int
	selectMax int
	waiting   int
)

func Init() error {
	state = stateBegin
	count = 0
	waiting = 0

	selectMax = 1
	if _, err := os.Stat(common.SaveFilePath); err == nil {
		selectMax = 2
	}

	cursor = selectMax - 1

	fname := common.ImagePath + "title/logo.png"
	imgLogo = dxlib.LoadGraph(fname)
	if imgBack == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	fname = common.ImagePath + "title/back.png"
	imgBack = dxlib.LoadGraph(fname)
	if imgBack == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}

	if err := sound.BGMPlay(sound.BGMTitle); err != nil {
		return fmt.Errorf("failed to play bgm: %v", err)
	}

	return nil
}

func End() {
	dxlib.DeleteGraph(imgLogo)
	dxlib.DeleteGraph(imgBack)
}

func Draw() {
	x := int32(-count % common.ScreenX)
	dxlib.DrawGraph(x, 0, imgBack, dxlib.FALSE)
	dxlib.DrawGraph(x+common.ScreenX, 0, imgBack, dxlib.FALSE)
	dxlib.DrawGraph(0, 0, imgLogo, dxlib.TRUE)

	switch state {
	case stateSelect:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenX, common.ScreenY, 0x000000, dxlib.TRUE)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)

		msgs := []string{"はじめから", "つづきから"}
		for i := 0; i < selectMax; i++ {
			draw.String(180, 230+int32(i)*20, 0xffffff, msgs[i])
		}
		const s = 2
		x := int32(160)
		y := int32(230 + cursor*20)
		dxlib.DrawTriangle(x, y+s, x+18-s*2, y+10, x, y+20-s, 0xffffff, dxlib.TRUE)
	}
}

func Process() error {
	if config.Get().Debug.SkipTitle {
		if config.Get().Debug.StartContinue {
			return ErrStartContinue
		}
		return ErrStartInit
	}

	switch state {
	case stateBegin:
		count++
		if count > 20 && inputs.CheckKey(inputs.KeyEnter) == 1 {
			state = stateSelect
		}
	case stateSelect:
		if waiting > 0 {
			waiting++
			if waiting > 30 {
				switch cursor {
				case 0:
					return ErrStartInit
				case 1:
					return ErrStartContinue
				default:
					return fmt.Errorf("unrecognized cursor %d was specified", cursor)
				}
			}
			return nil
		}

		if inputs.CheckKey(inputs.KeyEnter) == 1 {
			sound.On(sound.SETitleEnter)
			waiting++
		} else if inputs.CheckKey(inputs.KeyUp) == 1 && cursor > 0 {
			cursor--
		} else if inputs.CheckKey(inputs.KeyDown) == 1 && cursor < selectMax-1 {
			cursor++
		}
	}
	return nil
}
