package field

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

type Background struct {
	Images   []int
	Count    int
	LoadType int
}

const (
	BGType秋原町 int = iota
)

func (b *Background) Init(typ int) error {
	b.Count = 0
	b.LoadType = typ

	basePath := common.ImagePath + "battle/background/"

	switch typ {
	case BGType秋原町:
		b.Images = make([]int, 8)
		fname := basePath + "back_image_秋原町.png"
		if res := dxlib.LoadDivGraph(fname, 8, 2, 4, 64, 64, b.Images); res == -1 {
			return fmt.Errorf("failed to load image %s", fname)
		}
		return nil
	}

	return fmt.Errorf("invalid background type %d was specified", typ)
}

func (b *Background) End() {
	for _, img := range b.Images {
		dxlib.DeleteGraph(img)
	}
}

func (b *Background) Draw() {
	if len(b.Images) == 0 {
		return
	}

	switch b.LoadType {
	case BGType秋原町:
		dxlib.DrawBox(0, 0, common.ScreenSize.X, common.ScreenSize.Y, dxlib.GetColor(0, 0, 160), true)

		n := (b.Count / 50) % len(b.Images)
		i := 0
		spaceX := 140
		spaceY := 80
		for y := 60; y < common.ScreenSize.Y; y += spaceY {
			for x := 0; x < common.ScreenSize.X; x += spaceX {
				ofsX := (y / spaceY) * (spaceX / 2)
				n = (n + i) % len(b.Images)
				i++
				dxlib.DrawRotaGraph(x+ofsX, y, 1, 0, b.Images[n], true)
			}
		}
	}
}

func (b *Background) Process() {
	b.Count++
}
