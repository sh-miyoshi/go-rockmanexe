package effect

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
)

const (
	TypeNone int = iota
	TypeHitSmall
	TypeHitBig
	TypeExplode
	TypeCannonHit
	TypeSpreadHit
	TypeVulcanHit1
	TypeVulcanHit2

	typeMax
)

const (
	explodeDelay = 2
)

var (
	imgHitSmallEffect   []int32
	imgHitBigEffect     []int32
	imgExplodeEffect    []int32
	imgCannonHitEffect  []int32
	imgSpreadHitEffect  []int32
	imgVulcanHit1Effect []int32
	imgVulcanHit2Effect []int32
	sounds              [typeMax]sound.SEType
)

type effect struct {
	ID   string
	X    int
	Y    int
	Type int

	count  int
	images []int32
	delay  int
	ofsX   int32
	ofsY   int32
}

type noEffect struct{}

func Init() error {
	imgHitSmallEffect = make([]int32, 4)
	fname := common.ImagePath + "battle/effect/hit_small.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, imgHitSmallEffect); res == -1 {
		return fmt.Errorf("failed to load hit small effect image %s", fname)
	}
	imgHitBigEffect = make([]int32, 6)
	fname = common.ImagePath + "battle/effect/hit_big.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, imgHitBigEffect); res == -1 {
		return fmt.Errorf("failed to load hit big effect image %s", fname)
	}
	imgExplodeEffect = make([]int32, 16)
	fname = common.ImagePath + "battle/effect/explode.png"
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, imgExplodeEffect); res == -1 {
		return fmt.Errorf("failed to load explode effect image %s", fname)
	}
	imgCannonHitEffect = make([]int32, 7)
	fname = common.ImagePath + "battle/effect/cannon_hit.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, imgCannonHitEffect); res == -1 {
		return fmt.Errorf("failed to load cannon hit effect image %s", fname)
	}
	imgSpreadHitEffect = make([]int32, 6)
	fname = common.ImagePath + "battle/effect/spread_hit.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, imgSpreadHitEffect); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	tmp := make([]int32, 8)
	fname = common.ImagePath + "battle/effect/vulcan_hit.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 50, 58, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	imgVulcanHit1Effect = []int32{}
	imgVulcanHit2Effect = []int32{}
	for i := 0; i < 4; i++ {
		imgVulcanHit1Effect = append(imgVulcanHit1Effect, tmp[i])
		imgVulcanHit2Effect = append(imgVulcanHit2Effect, tmp[i+4])
	}

	for i := 0; i < typeMax; i++ {
		sounds[i] = -1
	}
	sounds[TypeCannonHit] = sound.SECannonHit
	sounds[TypeHitSmall] = sound.SEBusterHit
	sounds[TypeHitBig] = sound.SEBusterHit

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
	for _, img := range imgVulcanHit1Effect {
		dxlib.DeleteGraph(img)
	}
	for _, img := range imgVulcanHit2Effect {
		dxlib.DeleteGraph(img)
	}
}

func Get(typ int, x, y int, randRange int) anim.Anim {
	ofsX := 0
	ofsY := 0
	if randRange > 0 {
		ofsX = rand.Intn(2*randRange) - randRange
		ofsY = rand.Intn(2*randRange) - randRange
	}

	res := &effect{
		ID:    uuid.New().String(),
		Type:  typ,
		X:     x,
		Y:     y,
		delay: 1,
		ofsX:  int32(ofsX),
		ofsY:  int32(ofsY),
	}

	switch typ {
	case TypeNone:
		return &noEffect{}
	case TypeHitSmall:
		res.images = imgHitSmallEffect
	case TypeHitBig:
		res.images = imgHitBigEffect
	case TypeExplode:
		res.images = imgExplodeEffect
		res.delay = explodeDelay
	case TypeCannonHit:
		res.images = imgCannonHitEffect
	case TypeSpreadHit:
		res.images = imgSpreadHitEffect
	case TypeVulcanHit1:
		res.images = imgVulcanHit1Effect
		res.ofsY -= 30
	case TypeVulcanHit2:
		res.images = imgVulcanHit2Effect
		res.ofsY -= 10
	default:
		panic(fmt.Sprintf("Effect type %d is not implement yet.", typ))
	}

	return res
}

func (e *effect) Process() (bool, error) {
	e.count++

	if e.count == 1 {
		if sounds[e.Type] != -1 {
			sound.On(sounds[e.Type])
		}
	}

	return e.count >= len(e.images)*e.delay, nil
}

func (e *effect) Draw() {
	imgNo := -1
	if e.count < len(e.images)*e.delay {
		imgNo = e.count / e.delay
	}

	x, y := battlecommon.ViewPos(e.X, e.Y)
	dxlib.DrawRotaGraph(x+e.ofsX, y+e.ofsY+15, 1, 0, e.images[imgNo], dxlib.TRUE)
}

func (e *effect) DamageProc(dm *damage.Damage) bool {
	return false
}

func (e *effect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.ID,
		PosX:     e.X,
		PosY:     e.Y,
		AnimType: anim.TypeEffect,
		ObjType:  anim.ObjTypeNone,
	}
}

func (e *noEffect) Process() (bool, error) {
	// Nothing to do, so return finish immediately
	return true, nil
}

func (e *noEffect) Draw() {
}

func (e *noEffect) DamageProc(dm *damage.Damage) bool {
	return false
}

func (e *noEffect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    uuid.New().String(), // set dummy param
		AnimType: anim.TypeEffect,
		ObjType:  anim.ObjTypeNone,
	}
}
