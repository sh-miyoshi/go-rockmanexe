package background

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	TypeNone int = iota
	Type秋原町
)

type info struct {
	Images  []int
	Count   int
	BGColor uint
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

func Process() {
	bgInfo.Process()
}

func (i *info) Init(typ int) error {
	i.Count = 0

	basePath := common.ImagePath + "battle/background/"

	switch typ {
	case Type秋原町:
		i.BGColor = dxlib.GetColor(0, 0, 160)
		i.Images = make([]int, 8)
		fname := basePath + "back_image_秋原町.png"
		if res := dxlib.LoadDivGraph(fname, 8, 2, 4, 64, 64, i.Images); res == -1 {
			return fmt.Errorf("failed to load image %s", fname)
		}
		return nil
	}

	return fmt.Errorf("invalid background type %d was specified", typ)
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

	dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, i.BGColor, true)

	n := (i.Count / 50) % len(i.Images)
	index := 0
	spaceX := 140
	spaceY := 80
	for y := 60; y < common.ScreenSize.Y; y += spaceY {
		for x := 0; x < common.ScreenSize.X; x += spaceX {
			ofsX := (y / spaceY) * (spaceX / 2)
			n = (n + index) % len(i.Images)
			index++
			dxlib.DrawRotaGraph(x+ofsX, y, 1, 0, i.Images[n], true)
		}
	}
}

func (i *info) Process() {
	if len(i.Images) == 0 {
		return
	}

	i.Count++
}