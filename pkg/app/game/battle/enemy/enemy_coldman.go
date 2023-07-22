package enemy

import (
	"fmt"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
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
	pm     EnemyParam
	images [coldmanActTypeMax][]int
	count  int
	state  int
}

func (e *enemyColdman) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = coldmanActTypeStand

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
	case coldmanActTypeDamage:
		if e.count == 4*coldmanDelays[coldmanActTypeDamage] {
			// e.waitCount = 20
			e.state = coldmanActTypeStand
			// e.nextState = coldmanActTypeMove
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
