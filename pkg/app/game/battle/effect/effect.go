package effect

import (
	"fmt"
	"math/rand"

	"github.com/google/uuid"
	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
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
	TypeWaterBomb
	TypeBlock

	typeMax
)

const (
	explodeDelay   = 2
	waterBombDelay = 2
)

var (
	images [typeMax][]int32
	sounds [typeMax]sound.SEType
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
	images[TypeHitSmall] = make([]int32, 4)
	fname := common.ImagePath + "battle/effect/hit_small.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, images[TypeHitSmall]); res == -1 {
		return fmt.Errorf("failed to load hit small effect image %s", fname)
	}
	images[TypeHitBig] = make([]int32, 6)
	fname = common.ImagePath + "battle/effect/hit_big.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, images[TypeHitBig]); res == -1 {
		return fmt.Errorf("failed to load hit big effect image %s", fname)
	}
	images[TypeExplode] = make([]int32, 16)
	fname = common.ImagePath + "battle/effect/explode.png"
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, images[TypeExplode]); res == -1 {
		return fmt.Errorf("failed to load explode effect image %s", fname)
	}
	images[TypeCannonHit] = make([]int32, 7)
	fname = common.ImagePath + "battle/effect/cannon_hit.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, images[TypeCannonHit]); res == -1 {
		return fmt.Errorf("failed to load cannon hit effect image %s", fname)
	}
	images[TypeSpreadHit] = make([]int32, 6)
	fname = common.ImagePath + "battle/effect/spread_hit.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, images[TypeSpreadHit]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	tmp := make([]int32, 8)
	fname = common.ImagePath + "battle/effect/vulcan_hit.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 50, 58, tmp); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	images[TypeVulcanHit1] = []int32{}
	images[TypeVulcanHit2] = []int32{}
	for i := 0; i < 4; i++ {
		images[TypeVulcanHit1] = append(images[TypeVulcanHit1], tmp[i])
		images[TypeVulcanHit2] = append(images[TypeVulcanHit2], tmp[i+4])
	}
	images[TypeWaterBomb] = make([]int32, 12)
	fname = common.ImagePath + "battle/effect/water_bomb.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 112, 94, images[TypeWaterBomb]); res == -1 {
		return fmt.Errorf("failed to load image %s", fname)
	}
	for i := 7; i < 12; i++ {
		images[TypeWaterBomb][i] = images[TypeWaterBomb][6]
	}

	images[TypeBlock] = make([]int32, 4)
	fname = common.ImagePath + "battle/effect/block.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, images[TypeBlock]); res == -1 {
		return fmt.Errorf("failed to load block effect image %s", fname)
	}

	for i := 0; i < typeMax; i++ {
		sounds[i] = -1
	}
	sounds[TypeCannonHit] = sound.SECannonHit
	sounds[TypeHitSmall] = sound.SEBusterHit
	sounds[TypeHitBig] = sound.SEBusterHit
	sounds[TypeBlock] = sound.SEBlock

	return nil
}

func End() {
	for _, imgs := range images {
		for _, img := range imgs {
			dxlib.DeleteGraph(img)
		}
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
		ID:     uuid.New().String(),
		Type:   typ,
		X:      x,
		Y:      y,
		delay:  1,
		ofsX:   int32(ofsX),
		ofsY:   int32(ofsY),
		images: images[typ],
	}

	switch typ {
	case TypeNone:
		return &noEffect{}
	case TypeExplode:
		res.delay = explodeDelay
	case TypeVulcanHit1:
		res.ofsY -= 30
	case TypeVulcanHit2:
		res.ofsY -= 10
	case TypeWaterBomb:
		res.delay = waterBombDelay
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

func (e *effect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.ID,
		PosX:     e.X,
		PosY:     e.Y,
		AnimType: anim.AnimTypeEffect,
	}
}

func (e *noEffect) Process() (bool, error) {
	// Nothing to do, so return finish immediately
	return true, nil
}

func (e *noEffect) Draw() {
}

func (e *noEffect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    uuid.New().String(), // set dummy param
		AnimType: anim.AnimTypeEffect,
	}
}
