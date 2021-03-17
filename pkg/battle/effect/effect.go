package effect

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/dxlib"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

var (
	imgHitEffect []int32
)

type HitEffect struct {
	X int
	Y int

	count int
}

func Init() error {
	imgHitEffect = make([]int32, 6)
	fname := common.ImagePath + "battle/skill/バスター_hit.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, imgHitEffect); res == -1 {
		return fmt.Errorf("Failed to load hit effect image %s", fname)
	}

	return nil
}

func End() {
	for _, img := range imgHitEffect {
		dxlib.DeleteGraph(img)
	}
}

func (h *HitEffect) Process() (bool, error) {
	h.count++
	return h.count >= len(imgHitEffect), nil
}

func (h *HitEffect) Draw() {
	imgNo := -1
	if h.count < len(imgHitEffect) {
		imgNo = h.count
	}

	x, y := battlecommon.ViewPos(h.X, h.Y)
	x = x - 20 + rand.Int31n(40)
	y = y - 10 + rand.Int31n(20)

	dxlib.DrawRotaGraph(x, y, 1, 0, imgHitEffect[imgNo], dxlib.TRUE)
}
