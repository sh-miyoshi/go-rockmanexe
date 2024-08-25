package enemy

import (
	"github.com/cockroachdb/errors"
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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/math"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
)

const (
	forteActTypeStand = iota
	forteActTypeMove
	forteActTypeShooting
	forteActTypeHellsRolling
	forteActTypeDarkArmBlade
	forteActTypeDamage

	forteActTypeMax
)

var (
	forteDelays = [forteActTypeMax]int{1, 1, 1, 1, 1, 1}
)

type enemyForte struct {
	pm     EnemyParam
	images [forteActTypeMax][]int
	count  int
	state  int
}

func (e *enemyForte) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = forteActTypeStand

	// Load Images
	name, ext := GetStandImageFile(IDForte)

	fname := name + "_all" + ext
	tmp := make([]int, 45)
	if res := dxlib.LoadDivGraph(fname, 45, 9, 5, 136, 172, tmp); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
	}
	e.images[forteActTypeStand] = make([]int, 1)
	e.images[forteActTypeMove] = make([]int, 6)
	e.images[forteActTypeShooting] = make([]int, 9)
	e.images[forteActTypeHellsRolling] = make([]int, 9)
	e.images[forteActTypeDarkArmBlade] = make([]int, 3)
	e.images[forteActTypeDamage] = make([]int, 1)

	e.images[forteActTypeStand][0] = tmp[0]
	for i := 0; i < 6; i++ {
		e.images[forteActTypeMove][i] = tmp[i]
	}
	for i := 0; i < 9; i++ {
		e.images[forteActTypeShooting][i] = tmp[9+i]
		e.images[forteActTypeHellsRolling][i] = tmp[18+i]
	}
	for i := 0; i < 3; i++ {
		e.images[forteActTypeDarkArmBlade][i] = tmp[27+i]
	}
	e.images[forteActTypeDamage][0] = tmp[36]

	cleanup := []int{6, 7, 8}
	for i := 30; i < len(tmp); i++ {
		if i != 36 {
			cleanup = append(cleanup, i)
		}
	}

	for _, t := range cleanup {
		dxlib.DeleteGraph(t)
	}

	return nil
}

func (e *enemyForte) End() {
	// Delete Images
	for i := 0; i < forteActTypeMax; i++ {
		for j := 0; j < len(e.images[i]); j++ {
			dxlib.DeleteGraph(e.images[i][j])
		}
		e.images[i] = []int{}
	}
}

func (e *enemyForte) Process() (bool, error) {
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

	// WIP

	e.count++
	return false, nil
}

func (e *enemyForte) Draw() {
	if e.pm.InvincibleCount/5%2 != 0 {
		return
	}

	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()

	ofs := [forteActTypeMax]point.Point{
		{X: 0, Y: -20}, // Stand
		{X: 0, Y: 0},   // Move
		{X: 0, Y: 0},   // Shooting
		{X: 0, Y: 0},   // HellsRolling
		{X: 0, Y: 0},   // DarkArmBlade
		{X: 0, Y: 0},   // Damage
	}

	if e.state == forteActTypeStand {
		ofs[e.state].Y -= math.MountainIndex(e.count/10%5, 5)
	}

	dxlib.DrawRotaGraph(view.X+ofs[e.state].X, view.Y+ofs[e.state].Y, 1, 0, *img, true)

	drawParalysis(view.X+ofs[e.state].X, view.Y+ofs[e.state].Y, *img, e.pm.ParalyzedCount)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyForte) DamageProc(dm *damage.Damage) bool {
	if damageProc(dm, &e.pm) {
		if !dm.BigDamage {
			return true
		}

		e.state = forteActTypeDamage
		e.pm.InvincibleCount = battlecommon.PlayerDefaultInvincibleTime
		e.count = 0
		return true
	}

	return false
}

func (e *enemyForte) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID:    e.pm.ObjectID,
			Pos:      e.pm.Pos,
			DrawType: anim.DrawTypeObject,
		},
		HP: e.pm.HP,
	}
}

func (e *enemyForte) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyForte) MakeInvisible(count int) {
	e.pm.InvincibleCount = count
}

func (e *enemyForte) getCurrentImagePointer() *int {
	n := (e.count / forteDelays[e.state])
	if n >= len(e.images[e.state]) {
		n = len(e.images[e.state]) - 1
	}
	return &e.images[e.state][n]
}
