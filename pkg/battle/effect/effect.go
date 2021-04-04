package effect

import (
	"fmt"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/common"
)

const (
	TypeNone int = iota
	TypeHitSmall
	TypeHitBig
	TypeExplode
	TypeCannonHit
	TypeSpreadHit
)

const (
	explodeDelay = 2
)

var (
	imgHitSmallEffect  []int32
	imgHitBigEffect    []int32
	imgExplodeEffect   []int32
	imgCannonHitEffect []int32
	imgSpreadHitEffect []int32
)

type effect struct {
	X int
	Y int

	count  int
	images []int32
	delay  int
}

type noEffect struct{}

func Init() error {
	imgHitSmallEffect = make([]int32, 4)
	fname := common.ImagePath + "battle/effect/hit_small.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, imgHitSmallEffect); res == -1 {
		return fmt.Errorf("Failed to load hit small effect image %s", fname)
	}
	imgHitBigEffect = make([]int32, 6)
	fname = common.ImagePath + "battle/effect/hit_big.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, imgHitBigEffect); res == -1 {
		return fmt.Errorf("Failed to load hit big effect image %s", fname)
	}
	imgExplodeEffect = make([]int32, 16)
	fname = common.ImagePath + "battle/effect/explode.png"
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, imgExplodeEffect); res == -1 {
		return fmt.Errorf("Failed to load explode effect image %s", fname)
	}
	imgCannonHitEffect = make([]int32, 7)
	fname = common.ImagePath + "battle/effect/cannon_hit.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, imgCannonHitEffect); res == -1 {
		return fmt.Errorf("Failed to load cannon hit effect image %s", fname)
	}
	imgSpreadHitEffect = make([]int32, 6)
	fname = common.ImagePath + "battle/effect/spread_hit.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, imgSpreadHitEffect); res == -1 {
		return fmt.Errorf("Failed to load image %s", fname)
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
	for _, img := range imgSpreadHitEffect {
		dxlib.DeleteGraph(img)
	}
}

func Get(typ int, x, y int) anim.Anim {
	switch typ {
	case TypeNone:
		return &noEffect{}
	case TypeHitSmall:
		return &effect{X: x, Y: y, images: imgHitSmallEffect, delay: 1}
	case TypeHitBig:
		return &effect{X: x, Y: y, images: imgHitBigEffect, delay: 1}
	case TypeExplode:
		return &effect{X: x, Y: y, images: imgExplodeEffect, delay: explodeDelay}
	case TypeCannonHit:
		return &effect{X: x, Y: y, images: imgCannonHitEffect, delay: 1}
	case TypeSpreadHit:
		return &effect{X: x, Y: y, images: imgSpreadHitEffect, delay: 1}
	}

	panic(fmt.Sprintf("Effect type %d is not implement yet.", typ))
}

func (e *effect) Process() (bool, error) {
	e.count++
	return e.count >= len(e.images)*e.delay, nil
}

func (e *effect) Draw() {
	imgNo := -1
	if e.count < len(e.images)*e.delay {
		imgNo = e.count / e.delay
	}

	x, y := battlecommon.ViewPos(e.X, e.Y)
	dxlib.DrawRotaGraph(x, y+15, 1, 0, e.images[imgNo], dxlib.TRUE)
}

func (e *effect) DamageProc(dm *damage.Damage) {
}

func (e *effect) GetParam() anim.Param {
	return anim.Param{
		PosX:     e.X,
		PosY:     e.Y,
		AnimType: anim.TypeEffect,
	}
}

func (e *noEffect) Process() (bool, error) {
	// Nothing to do, so return finish immediately
	return true, nil
}

func (e *noEffect) Draw() {
}

func (e *noEffect) DamageProc(dm *damage.Damage) {
}

func (e *noEffect) GetParam() anim.Param {
	return anim.Param{
		AnimType: anim.TypeEffect,
	}
}
