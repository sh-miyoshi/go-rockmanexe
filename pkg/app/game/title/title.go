package title

import (
	"errors"
	"fmt"
	"os"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/inputs"
)

const (
	stateBegin int = iota
	stateSelect
)

var (
	ErrStartInit     = errors.New("start with initialize")
	ErrStartContinue = errors.New("start with continue")

	count     int
	imgLogo   int
	imgBack   int
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
	x := -count % common.ScreenSize.X
	dxlib.DrawGraph(x, 0, imgBack, false)
	dxlib.DrawGraph(x+common.ScreenSize.X, 0, imgBack, false)
	dxlib.DrawGraph(0, 0, imgLogo, true)

	switch state {
	case stateSelect:
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_ALPHA, 192)
		dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, 0x000000, true)
		dxlib.SetDrawBlendMode(dxlib.DX_BLENDMODE_NOBLEND, 255)

		msgs := []string{"はじめから", "つづきから"}
		for i := 0; i < selectMax; i++ {
			draw.String(180, 230+i*20, 0xffffff, msgs[i])
		}
		const s = 2
		x := 160
		y := 230 + cursor*20
		dxlib.DrawTriangle(x, y+s, x+18-s*2, y+10, x, y+20-s, 0xffffff, true)
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
