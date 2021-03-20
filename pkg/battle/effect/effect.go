package effect

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

const (
	TypeNone int = iota
	TypeHitSmall
	TypeHitBig
	TypeExplode
	TypeCannonHit
)

const (
	explodeDelay = 2
)

var (
	imgHitSmallEffect  []int32
	imgHitBigEffect    []int32
	imgExplodeEffect   []int32
	imgCannonHitEffect []int32
)

type effect struct {
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
	imgExplodeEffect = make([]int32, 16)
	fname = common.ImagePath + "battle/skill/explode.png"
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, imgExplodeEffect); res == -1 {
		return fmt.Errorf("Failed to load explode effect image %s", fname)
	}
	imgCannonHitEffect = make([]int32, 7)
	fname = common.ImagePath + "battle/skill/キャノン_hit.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, imgCannonHitEffect); res == -1 {
		return fmt.Errorf("Failed to load cannon hit effect image %s", fname)
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
	for _, img := range imgExplodeEffect {
		dxlib.DeleteGraph(img)
	}
	for _, img := range imgCannonHitEffect {
		dxlib.DeleteGraph(img)
	}
}

func Get(typ int, x, y int) anim.Anim {
	switch typ {
	case TypeHitSmall:
		return &effect{X: x, Y: y, images: imgHitSmallEffect}
	case TypeHitBig:
		return &effect{X: x, Y: y, images: imgHitBigEffect}
	case TypeExplode:
		return &effect{X: x, Y: y, images: imgExplodeEffect}
	case TypeCannonHit:
		return &effect{X: x, Y: y, images: imgCannonHitEffect}
	}
	return nil
}

func (e *effect) Process() (bool, error) {
	e.count++
	return e.count >= len(e.images)*explodeDelay, nil
}

func (e *effect) Draw() {
	imgNo := -1
	if e.count < len(e.images)*explodeDelay {
		imgNo = e.count / explodeDelay
	}

	x, y := battlecommon.ViewPos(e.X, e.Y)
	dxlib.DrawRotaGraph(x, y+15, 1, 0, e.images[imgNo], dxlib.TRUE)
}
