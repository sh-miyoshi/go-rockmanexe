package field

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type background struct {
	Images   []int
	Count    int
	LoadType int
}

const (
	bgType秋原町 int = iota
)

func (b *background) Init(typ int) error {
	b.Count = 0
	b.LoadType = typ

	basePath := common.ImagePath + "battle/background/"

	switch typ {
	case bgType秋原町:
		b.Images = make([]int, 8)
		fname := basePath + "back_image_秋原町.png"
		if res := dxlib.LoadDivGraph(fname, 8, 2, 4, 64, 64, b.Images); res == -1 {
			return fmt.Errorf("failed to load image %s", fname)
		}
		return nil
	}

	return fmt.Errorf("invalid background type %d was specified", typ)
}

func (b *background) End() {
	for _, img := range b.Images {
		dxlib.DeleteGraph(img)
	}
}

func (b *background) Draw() {
	if len(b.Images) == 0 {
		return
	}

	switch b.LoadType {
	case bgType秋原町:
		dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, dxlib.GetColor(0, 0, 160), true)

		n := (b.Count / 50) % len(b.Images)
		initPos := []common.Point{
			{X: 0, Y: 60},
			{X: 140, Y: 60},
			{X: 280, Y: 60},
			{X: 420, Y: 60},
			{X: 30, Y: 140},
			{X: 170, Y: 140},
			{X: 310, Y: 140},
			{X: 450, Y: 140},
		}

		for i, pos := range initPos {
			n = (n + i) % len(b.Images)
			dxlib.DrawRotaGraph(pos.X, pos.Y, 1, 0, b.Images[n], true)
		}
	}
}

func (b *background) Process() {
	b.Count++
}
