package background

import (
	"github.com/cockroachdb/errors"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	TypeNone int = iota
	Type秋原町
	Typeアッフリク
	Typeブラックアース
)

type info struct {
	Images  []int
	Count   int
	BGColor uint
	Type    int
}

var bgInfo info

func Set(bgType int) error {
	bgInfo.End()
	return bgInfo.Init(bgType)
}

func Unset() {
	bgInfo.End()
}

func Draw() {
	bgInfo.Draw()
}

func Update() {
	bgInfo.Update()
}

func (i *info) Init(typ int) error {
	i.End()

	i.Count = 0
	i.Type = typ

	basePath := config.ImagePath + "battle/background/"

	switch typ {
	case Type秋原町:
		i.BGColor = dxlib.GetColor(0, 0, 160)
		i.Images = make([]int, 8)
		fname := basePath + "back_image_秋原町.png"
		if res := dxlib.LoadDivGraph(fname, 8, 2, 4, 64, 64, i.Images); res == -1 {
			return errors.Newf("failed to load image %s", fname)
		}
		return nil
	case Typeアッフリク:
		i.BGColor = dxlib.GetColor(255, 140, 0)
		i.Images = make([]int, 8)
		fname := basePath + "back_image_アッフリク.png"
		if res := dxlib.LoadDivGraph(fname, 8, 2, 4, 64, 64, i.Images); res == -1 {
			return errors.Newf("failed to load image %s", fname)
		}
		return nil
	case Typeブラックアース:
		i.BGColor = 0
		i.Images = make([]int, 1)
		fname := basePath + "back_image_ブラックアース.png"
		i.Images[0] = dxlib.LoadGraph(fname)
		if i.Images[0] == -1 {
			return errors.Newf("failed to load image %s", fname)
		}
		return nil
	}

	return errors.Newf("invalid background type %d was specified", typ)
}

func (i *info) End() {
	for _, img := range i.Images {
		dxlib.DeleteGraph(img)
	}
	i.Images = []int{}
}

func (i *info) Draw() {
	if len(i.Images) == 0 {
		return
	}

	switch i.Type {
	case Type秋原町, Typeアッフリク:
		dxlib.DrawBox(0, 0, config.ScreenSize.X, config.ScreenSize.Y, i.BGColor, true)

		n := (i.Count / 50) % len(i.Images)
		index := 0
		spaceX := 140
		spaceY := 80
		for y := 60; y < config.ScreenSize.Y; y += spaceY {
			for x := 0; x < config.ScreenSize.X; x += spaceX {
				ofsX := (y / spaceY) * (spaceX / 2)
				n = (n + index) % len(i.Images)
				index++
				dxlib.DrawRotaGraph(x+ofsX, y, 1, 0, i.Images[n], true)
			}
		}
	case Typeブラックアース:
		scroll := (i.Count / 10) % config.ScreenSize.X
		dxlib.DrawGraph(scroll-config.ScreenSize.X, 0, i.Images[0], false)
		dxlib.DrawRotaGraph(scroll, 0, 1, 0, i.Images[0], false, dxlib.OptXReverse(true))
	}
}

func (i *info) Update() {
	if len(i.Images) == 0 {
		return
	}

	i.Count++
}
