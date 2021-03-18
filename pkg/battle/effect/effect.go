package effect

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

const (
	TypeNone int = iota
	TypeHitSmall
	TypeHitBig
)

var (
	imgHitSmallEffect []int32
	imgHitBigEffect   []int32
)

type HitEffect struct {
	X int
	Y int

	count  int
	images []int32
}

func Init() error {
	imgHitSmallEffect = make([]int32, 4)
	fname := common.ImagePath + "battle/skill/バスター_hit_small.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, imgHitSmallEffect); res == -1 {
		return fmt.Errorf("Failed to load hit small effect image %s", fname)
	}
	imgHitBigEffect = make([]int32, 6)
	fname = common.ImagePath + "battle/skill/バスター_hit_big.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, imgHitBigEffect); res == -1 {
		return fmt.Errorf("Failed to load hit big effect image %s", fname)
	}

	return nil
}

func End() {
	for _, img := range imgHitSmallEffect {
		dxlib.DeleteGraph(img)
	}
	for _, img := range imgHitBigEffect {
		dxlib.DeleteGraph(img)
	}
}

func Get(typ int, x, y int) anim.Anim {
	switch typ {
	case TypeHitSmall:
		return &HitEffect{X: x, Y: y, images: imgHitSmallEffect}
	case TypeHitBig:
		return &HitEffect{X: x, Y: y, images: imgHitBigEffect}
	}
	return nil
}

func (h *HitEffect) Process() (bool, error) {
	h.count++
	return h.count >= len(h.images), nil
}

func (h *HitEffect) Draw() {
	imgNo := -1
	if h.count < len(h.images) {
		imgNo = h.count
	}

	x, y := battlecommon.ViewPos(h.X, h.Y)
	x = x - 15 + rand.Int31n(30)
	y = y - 5 + rand.Int31n(10)

	dxlib.DrawRotaGraph(x, y, 1, 0, h.images[imgNo], dxlib.TRUE)
}
