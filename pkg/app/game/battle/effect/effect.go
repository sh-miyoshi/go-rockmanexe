package effect

import (
	"math/rand"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/config"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/sound"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	explodeDelay     = 2
	waterBombDelay   = 2
	exclamationDelay = 2
)

var (
	images [resources.EffectTypeMax][]int
	sounds [resources.EffectTypeMax]resources.SEType
)

type effect struct {
	ID   string
	Pos  point.Point
	Type int

	count  int
	images []int
	delay  int
	ofs    point.Point
}

type noEffect struct{}

func Init() error {
	images[resources.EffectTypeHitSmall] = make([]int, 4)
	fname := config.ImagePath + "battle/effect/hit_small.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, images[resources.EffectTypeHitSmall]); res == -1 {
		return errors.Newf("failed to load hit small effect image %s", fname)
	}
	images[resources.EffectTypeHitBig] = make([]int, 6)
	fname = config.ImagePath + "battle/effect/hit_big.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 90, 76, images[resources.EffectTypeHitBig]); res == -1 {
		return errors.Newf("failed to load hit big effect image %s", fname)
	}
	images[resources.EffectTypeExplode] = make([]int, 16)
	fname = config.ImagePath + "battle/effect/explode.png"
	if res := dxlib.LoadDivGraph(fname, 16, 8, 2, 110, 124, images[resources.EffectTypeExplode]); res == -1 {
		return errors.Newf("failed to load explode effect image %s", fname)
	}
	images[resources.EffectTypeCannonHit] = make([]int, 7)
	fname = config.ImagePath + "battle/effect/cannon_hit.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 110, 136, images[resources.EffectTypeCannonHit]); res == -1 {
		return errors.Newf("failed to load cannon hit effect image %s", fname)
	}
	images[resources.EffectTypeSpreadHit] = make([]int, 6)
	fname = config.ImagePath + "battle/effect/spread_and_bamboo_hit.png"
	if res := dxlib.LoadDivGraph(fname, 6, 6, 1, 92, 88, images[resources.EffectTypeSpreadHit]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	tmp := make([]int, 8)
	fname = config.ImagePath + "battle/effect/vulcan_hit.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 50, 58, tmp); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	images[resources.EffectTypeVulcanHit1] = []int{}
	images[resources.EffectTypeVulcanHit2] = []int{}
	for i := 0; i < 4; i++ {
		images[resources.EffectTypeVulcanHit1] = append(images[resources.EffectTypeVulcanHit1], tmp[i])
		images[resources.EffectTypeVulcanHit2] = append(images[resources.EffectTypeVulcanHit2], tmp[i+4])
	}
	images[resources.EffectTypeWaterBomb] = make([]int, 12)
	fname = config.ImagePath + "battle/effect/water_bomb.png"
	if res := dxlib.LoadDivGraph(fname, 7, 7, 1, 112, 94, images[resources.EffectTypeWaterBomb]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	for i := 7; i < 12; i++ {
		images[resources.EffectTypeWaterBomb][i] = images[resources.EffectTypeWaterBomb][6]
	}

	images[resources.EffectTypeBlock] = make([]int, 4)
	fname = config.ImagePath + "battle/effect/block.png"
	if res := dxlib.LoadDivGraph(fname, 4, 4, 1, 40, 44, images[resources.EffectTypeBlock]); res == -1 {
		return errors.Newf("failed to load block effect image %s", fname)
	}

	images[resources.EffectTypeBambooHit] = append([]int{}, images[resources.EffectTypeSpreadHit]...)
	images[resources.EffectTypeHeatHit] = append([]int{}, images[resources.EffectTypeCannonHit]...)
	images[resources.EffectTypeExclamation] = make([]int, 6)
	fname = config.ImagePath + "battle/effect/exclamation.png"
	if res := dxlib.LoadDivGraph(fname, 3, 3, 1, 104, 102, images[resources.EffectTypeExclamation]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	for i := 3; i < 6; i++ {
		images[resources.EffectTypeExclamation][i] = images[resources.EffectTypeExclamation][2]
	}
	images[resources.EffectTypeFailed] = make([]int, 8)
	fname = config.ImagePath + "battle/effect/failed.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 38, 38, images[resources.EffectTypeFailed]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	images[resources.EffectTypeIceBreak] = make([]int, 2)
	fname = config.ImagePath + "battle/effect/ice_break.png"
	if res := dxlib.LoadDivGraph(fname, 2, 2, 1, 50, 56, images[resources.EffectTypeIceBreak]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}
	images[resources.EffectTypeExplodeSmall] = make([]int, 8)
	fname = config.ImagePath + "battle/effect/explode_small.png"
	if res := dxlib.LoadDivGraph(fname, 8, 8, 1, 84, 80, images[resources.EffectTypeExplodeSmall]); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	for i := 0; i < resources.EffectTypeMax; i++ {
		sounds[i] = -1
	}
	sounds[resources.EffectTypeCannonHit] = resources.SECannonHit
	sounds[resources.EffectTypeHitSmall] = resources.SEBusterHit
	sounds[resources.EffectTypeHitBig] = resources.SEBusterHit
	sounds[resources.EffectTypeBlock] = resources.SEBlock
	sounds[resources.EffectTypeHeatHit] = resources.SEExplode
	// TODO add exclamation se

	return nil
}

func End() {
	for _, imgs := range images {
		for _, img := range imgs {
			dxlib.DeleteGraph(img)
		}
	}
}

func Get(typ int, pos point.Point, randRange int) anim.Anim {
	ofs := point.Point{}
	if randRange > 0 {
		ofs.X = rand.Intn(2*randRange) - randRange
		ofs.Y = rand.Intn(2*randRange) - randRange
	}

	res := &effect{
		ID:     uuid.New().String(),
		Type:   typ,
		Pos:    pos,
		delay:  1,
		ofs:    ofs,
		images: images[typ],
	}

	switch typ {
	case resources.EffectTypeNone:
		return &noEffect{}
	case resources.EffectTypeExplode:
		res.delay = explodeDelay
	case resources.EffectTypeVulcanHit1:
		res.ofs.Y -= 30
	case resources.EffectTypeVulcanHit2:
		res.ofs.Y -= 10
	case resources.EffectTypeWaterBomb:
		res.delay = waterBombDelay
	case resources.EffectTypeExclamation:
		res.delay = exclamationDelay
		res.ofs.Y -= 40
	case resources.EffectTypeFailed:
		res.ofs.Y -= 60
	case resources.EffectTypeIceBreak:
		return &iceBreakEffect{
			ID:  uuid.New().String(),
			Pos: pos,
			Ofs: ofs,
		}
	case resources.EffectTypeExplodeSmall:
		system.SetError("explode small effect is not implemented yet")
	case resources.EffectTypeSpecialStart:
		res := &specialStartEffect{
			ID:  uuid.New().String(),
			Pos: pos,
			Ofs: ofs,
		}
		res.Init()
		return res
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

	view := battlecommon.ViewPos(e.Pos)
	dxlib.DrawRotaGraph(view.X+e.ofs.X, view.Y+e.ofs.Y+15, 1, 0, e.images[imgNo], true)
}

func (e *effect) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.ID,
		Pos:      e.Pos,
		DrawType: anim.DrawTypeEffect,
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
		DrawType: anim.DrawTypeEffect,
	}
}
