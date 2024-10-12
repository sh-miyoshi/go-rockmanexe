package enemy

import (
	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	bluesActTypeStand = iota
	bluesActTypeMove
	bluesActTypeWideSword
	bluesActTypeFighterSword
	bluesActTypeSonicBoom
	bluesActTypeDeltaRayEdge
	bluesActTypeBehindSlash
	bluesActTypeDamage

	bluesActTypeMax
)

var (
	bluesDelays = [bluesActTypeMax]int{1, 2, 1, 1, 1, 1, 1}
)

type enemyBlues struct {
	pm        EnemyParam
	state     int
	count     int
	waitCount int
	nextState int
	targetPos point.Point
	images    [bluesActTypeMax][]int
}

func (e *enemyBlues) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = bluesActTypeStand
	e.count = 0
	e.waitCount = 20
	e.nextState = bluesActTypeMove

	// Load Images
	name, ext := GetStandImageFile(IDBlues)

	fname := name + "_all" + ext
	tmp := make([]int, 36)
	if res := dxlib.LoadDivGraph(fname, 36, 7, 6, 170, 156, tmp); res == -1 {
		return errors.Newf("failed to load image %s", fname)
	}

	releases := [36]int{}
	for i := 0; i < 36; i++ {
		releases[i] = i
	}

	e.images[bluesActTypeStand] = make([]int, 1)
	e.images[bluesActTypeStand][0] = tmp[0]
	releases[0] = -1

	e.images[bluesActTypeMove] = make([]int, 4)
	for i := 0; i < 4; i++ {
		e.images[bluesActTypeMove][i] = tmp[i]
		releases[i] = -1
	}

	e.images[bluesActTypeWideSword] = make([]int, 6)
	e.images[bluesActTypeFighterSword] = make([]int, 6)
	for i := 0; i < 6; i++ {
		e.images[bluesActTypeWideSword][i] = tmp[i+7]
		e.images[bluesActTypeFighterSword][i] = tmp[i+7]
		releases[i+7] = -1
	}

	// 使わないイメージを削除
	for i, r := range releases {
		if r != -1 {
			dxlib.DeleteGraph(tmp[i])
		}
	}

	return nil
}

func (e *enemyBlues) End() {
	// Delete Images
	for _, imgs := range e.images {
		for _, img := range imgs {
			dxlib.DeleteGraph(img)
		}
	}
}

func (e *enemyBlues) Process() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false)
		localanim.AnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	if e.pm.ParalyzedCount > 0 {
		e.pm.ParalyzedCount--
		return false, nil
	}

	// Enemy Logic
	if e.pm.InvincibleCount > 0 {
		e.pm.InvincibleCount--
	}

	switch e.state {
	case bluesActTypeStand:
		e.waitCount--
		if e.waitCount <= 0 {
			return e.stateChange(e.nextState)
		}
	case bluesActTypeMove:
		if e.count == 3*bluesDelays[bluesActTypeMove] {
			// WIP: targetPos
			moveRandom(&e.pm.Pos)
			e.waitCount = 40

			// WIP: attack

			return e.stateChange(bluesActTypeStand)
		}
	}

	e.count++
	return false, nil
}

func (e *enemyBlues) Draw() {
	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()
	ofs := [bluesActTypeMax]point.Point{
		{X: -5, Y: -20}, // Stand
		{X: -5, Y: -20}, // Move
		{X: 0, Y: 0},    // Sword
		{X: 0, Y: 0},    // Shot
		{X: 0, Y: 0},    // Throw
		{X: 0, Y: 0},    // Throw2
		{X: 0, Y: 0},    // Damage
	}

	flag := int32(dxlib.TRUE)
	opt := dxlib.DrawRotaGraphOption{
		ReverseXFlag: &flag,
	}

	dxlib.DrawRotaGraph(view.X+ofs[e.state].X, view.Y+ofs[e.state].Y, 1, 0, *img, true, opt)
}

func (e *enemyBlues) DamageProc(dm *damage.Damage) bool {
	return damageProc(dm, &e.pm)
}

func (e *enemyBlues) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyBlues) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyBlues) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyBlues) getCurrentImagePointer() *int {
	if e.count == 0 {
		return &e.images[bluesActTypeStand][0]
	}

	n := (e.count / bluesDelays[e.state])
	if n >= len(e.images[e.state]) {
		n = len(e.images[e.state]) - 1
	}
	return &e.images[e.state][n]
}

func (e *enemyBlues) stateChange(next int) (bool, error) {
	logger.Info("change blues state to %d", next)
	e.state = next
	e.count = 0

	return false, nil
}
