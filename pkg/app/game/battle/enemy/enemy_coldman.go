package enemy

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
)

const (
	coldmanActTypeStand = iota
	coldmanActTypeIceCreate
	coldmanActTypeMove
	coldmanActTypeIceShoot
	coldmanActTypeBodyBlow
	coldmanActTypeBless
	coldmanActTypeDamage

	coldmanActTypeMax
)

var (
	coldmanDelays = [coldmanActTypeMax]int{1, 1, 2, 1, 1, 1, 5}
)

type enemyColdman struct {
	pm        EnemyParam
	images    [coldmanActTypeMax][]int
	count     int
	state     int
	nextState int
	waitCount int
	moveNum   int
	cubeIDs   []string
}

func (e *enemyColdman) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = coldmanActTypeStand
	e.waitCount = 60
	e.nextState = coldmanActTypeMove
	e.moveNum = rand.Intn(2) + 2
	e.cubeIDs = []string{}

	// Load Images
	name, ext := GetStandImageFile(IDColdman)

	fname := name + "_all" + ext
	tmp := make([]int, 24)
	if res := dxlib.LoadDivGraph(fname, 24, 6, 4, 136, 115, tmp); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}
	cleanup := []int{}
	e.images[coldmanActTypeStand] = make([]int, 1)
	e.images[coldmanActTypeStand][0] = tmp[0]
	e.images[coldmanActTypeIceCreate] = make([]int, 1)
	e.images[coldmanActTypeIceCreate][0] = tmp[0]

	e.images[coldmanActTypeMove] = make([]int, 2)
	e.images[coldmanActTypeIceShoot] = make([]int, 4)
	e.images[coldmanActTypeBodyBlow] = make([]int, 6)
	for j := 0; j < 3; j++ {
		for i := 0; i < 6; i++ {
			if i < len(e.images[j+coldmanActTypeMove]) {
				e.images[j+coldmanActTypeMove][i] = tmp[j*6+i]
			} else {
				cleanup = append(cleanup, j*6+i)
			}
		}
	}

	e.images[coldmanActTypeBless] = make([]int, 3)
	for i := 0; i < 3; i++ {
		e.images[coldmanActTypeBless][i] = tmp[18+i]
	}
	e.images[coldmanActTypeDamage] = make([]int, 1)
	e.images[coldmanActTypeDamage][0] = tmp[21]
	for i := 21; i < 24; i++ {
		cleanup = append(cleanup, i)
	}

	for _, t := range cleanup {
		dxlib.DeleteGraph(t)
	}

	return nil
}

func (e *enemyColdman) End() {
	// Delete Images
	for i := 0; i < coldmanActTypeMax; i++ {
		for j := 0; j < len(e.images[i]); j++ {
			dxlib.DeleteGraph(e.images[i][j])
		}
		e.images[i] = []int{}
	}
}

func (e *enemyColdman) Process() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false)
		localanim.AnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

	// Enemy Logic
	if e.pm.InvincibleCount > 0 {
		e.pm.InvincibleCount--
	}

	switch e.state {
	case coldmanActTypeStand:
		e.waitCount--
		if e.waitCount <= 0 {
			e.state = e.nextState
			e.count = 0
			return false, nil
		}
	case coldmanActTypeMove:
		if e.count == 2*coldmanDelays[coldmanActTypeMove] {
			e.moveRandom()

			e.waitCount = 60
			e.state = coldmanActTypeStand
			e.moveNum--
			if e.moveNum <= 0 {
				e.moveNum = rand.Intn(2) + 2

				// TODO next action
				e.nextState = coldmanActTypeIceCreate
			}
		}
	case coldmanActTypeIceCreate:
		if e.count == 0 {
			field.SetBlackoutCount(90)

			if err := e.createCube(); err != nil {
				return false, nil
			}
		}
	case coldmanActTypeIceShoot:
		panic("not implemented yet")
	case coldmanActTypeBodyBlow:
		panic("not implemented yet")
	case coldmanActTypeBless:
		panic("not implemented yet")
	case coldmanActTypeDamage:
		if e.count == 4*coldmanDelays[coldmanActTypeDamage] {
			e.waitCount = 20
			e.state = coldmanActTypeStand
			e.nextState = coldmanActTypeMove
			e.count = 0
			return false, nil
		}
	}

	e.count++
	return false, nil
}

func (e *enemyColdman) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()

	ofs := [coldmanActTypeMax]common.Point{
		{X: 0, Y: 0},  // Stand
		{X: 0, Y: 0},  // IceCreate
		{X: 0, Y: 0},  // Move
		{X: 0, Y: 0},  // IceShoot
		{X: 0, Y: 0},  // BodyBlow
		{X: 0, Y: 0},  // Bless
		{X: 20, Y: 0}, // Damage
	}

	dxlib.DrawRotaGraph(view.X+ofs[e.state].X, view.Y+ofs[e.state].Y, 1, 0, *img, true)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyColdman) DamageProc(dm *damage.Damage) bool {
	if damageProc(dm, &e.pm) {
		if !dm.BigDamage {
			return true
		}

		e.state = coldmanActTypeDamage
		e.pm.InvincibleCount = battlecommon.PlayerDefaultInvincibleTime
		e.count = 0
		return true
	}

	return false
}

func (e *enemyColdman) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyColdman) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyColdman) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyColdman) getCurrentImagePointer() *int {
	n := (e.count / coldmanDelays[e.state])
	if n >= len(e.images[e.state]) {
		n = len(e.images[e.state]) - 1
	}
	return &e.images[e.state][n]
}

func (e *enemyColdman) moveRandom() {
	// 移動先は最後列のどこか
	x := battlecommon.FieldNum.X - 1
	for i := 0; i < 10; i++ {
		next := common.Point{
			X: x,
			Y: rand.Intn(battlecommon.FieldNum.Y),
		}
		if battlecommon.MoveObjectDirect(
			&e.pm.Pos,
			next,
			battlecommon.PanelTypeEnemy,
			true,
			field.GetPanelInfo,
		) {
			return
		}
	}
}

func (e *enemyColdman) createCube() error {
	pm := object.ObjectParam{
		Pos:           common.Point{X: 4, Y: 1}, // TODO(特定のパターンで3個生成)
		HP:            200,
		OnwerCharType: objanim.ObjTypeEnemy,
	}
	obj := &object.IceCube{}
	if err := obj.Init(e.pm.ObjectID, pm); err != nil {
		return fmt.Errorf("failed to init ice cube: %w", err)
	}
	id := localanim.ObjAnimNew(obj)
	localanim.ObjAnimAddActiveAnim(id)

	e.cubeIDs = append(e.cubeIDs, id)

	return nil
}
