package enemy

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
)

const (
	aquamanActTypeStand = iota
	aquamanActTypeMove
	aquamanActTypeShot
	aquamanActTypeDamage
	aquamanActTypeBomb

	aquamanActTypeMax
)

var (
	aquamanDelays = [aquamanActTypeMax]int{8, 4, 4, 4, 3}
)

type enemyAquaman struct {
	pm              EnemyParam
	images          [aquamanActTypeMax][]int32
	count           int
	state           int
	nextState       int
	waitCount       int
	invincibleCount int
}

func (e *enemyAquaman) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = aquamanActTypeStand
	e.waitCount = 10
	e.nextState = aquamanActTypeBomb

	// Load Images
	name, ext := GetStandImageFile(IDAquaman)

	fname := name + "_stand" + ext
	e.images[aquamanActTypeStand] = make([]int32, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, e.images[aquamanActTypeStand]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_move" + ext
	e.images[aquamanActTypeMove] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 64, 126, e.images[aquamanActTypeMove]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_shot" + ext
	e.images[aquamanActTypeShot] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 104, 90, e.images[aquamanActTypeShot]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_damage" + ext
	e.images[aquamanActTypeDamage] = make([]int32, 1)
	if res := dxlib.LoadDivGraph(fname, 1, 1, 1, 70, 86, e.images[aquamanActTypeDamage]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_bomb" + ext
	e.images[aquamanActTypeBomb] = make([]int32, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 124, e.images[aquamanActTypeBomb]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	return nil
}

func (e *enemyAquaman) End() {
	// Delete Images
	for i := 0; i < aquamanActTypeMax; i++ {
		for j := 0; j < len(e.images[i]); j++ {
			dxlib.DeleteGraph(e.images[i][j])
		}
		e.images[i] = []int32{}
	}
}

func (e *enemyAquaman) Process() (bool, error) {
	// Return true if finished(e.g. hp=0)
	// Enemy Logic
	if e.invincibleCount > 0 {
		e.invincibleCount++
		if e.invincibleCount > battlecommon.PlayerDefaultInvincibleTime {
			e.invincibleCount = 0
		}
	}

	switch e.state {
	case aquamanActTypeStand:
		e.waitCount--
		if e.waitCount <= 0 {
			e.state = e.nextState
			e.count = 0
			return false, nil
		}
	case aquamanActTypeMove:
		if e.count == 5*aquamanDelays[aquamanActTypeMove] {
			for i := 0; i < 10; i++ {
				if battlecommon.MoveObjectDirect(
					&e.pm.PosX,
					&e.pm.PosY,
					rand.Intn(field.FieldNumX/2)+field.FieldNumX/2,
					rand.Intn(field.FieldNumY),
					field.PanelTypeEnemy,
					true,
					field.GetPanelInfo,
				) {
					break
				}
			}
			e.waitCount = 60
			e.state = aquamanActTypeStand
			e.count = 0
			return false, nil
		}
	case aquamanActTypeShot:
		// TODO(sound)

		if e.count == 6*aquamanDelays[aquamanActTypeShot] {
			for i := 0; i < 2; i++ {
				x := e.pm.PosX - (2 + i)
				anim.New(effect.Get(effect.TypeWaterBomb, x, e.pm.PosY, 0))
				damage.New(damage.Damage{
					PosX:          x,
					PosY:          e.pm.PosY,
					Power:         20, // TODO(ダメージ量)
					TTL:           1,
					TargetType:    damage.TargetPlayer,
					HitEffectType: effect.TypeNone,
					BigDamage:     true,
				})
			}

			e.waitCount = 60
			e.state = aquamanActTypeStand
			e.count = 0
			return false, nil
		}
	case aquamanActTypeDamage:
		if e.count == 4*aquamanDelays[aquamanActTypeDamage] {
			e.waitCount = 20
			e.state = aquamanActTypeStand
			e.nextState = aquamanActTypeMove // TODO(次の行動)
			e.count = 0
			return false, nil
		}
	case aquamanActTypeBomb:
		if e.count == 3*aquamanDelays[aquamanActTypeBomb] {
			// ボム登録
			anim.New(skill.Get(skill.SkillWaterBomb, skill.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      20, // TODO(ダメージ量)
				TargetType: damage.TargetPlayer,
			}))
		}

		if e.count == 6*aquamanDelays[aquamanActTypeBomb] {
			e.waitCount = 120
			e.state = aquamanActTypeStand
			e.count = 0
			// TODO(次の行動)
			return false, nil
		}
	}

	e.count++
	return false, nil
}

func (e *enemyAquaman) Draw() {
	if e.invincibleCount/5%2 != 0 {
		return
	}

	// Show Enemy Images
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.getCurrentImagePointer()

	ofs := [aquamanActTypeMax][]int32{
		{0, 0},    // stand
		{0, 0},    // move
		{-20, 10}, // shot
		{0, 0},    // damage
		{0, 0},    // bomb
	}

	dxlib.DrawRotaGraph(x+ofs[e.state][0], y+ofs[e.state][1], 1, 0, *img, dxlib.TRUE)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(x, y+40, int32(e.pm.HP), draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyAquaman) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}
	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.PosX, e.pm.PosY, 5))

		if !dm.BigDamage {
			return true
		}

		e.state = aquamanActTypeDamage
		e.invincibleCount = 1
		e.count = 0
		return true
	}
	return false
}

func (e *enemyAquaman) GetParam() anim.Param {
	return anim.Param{
		ObjID:    e.pm.ObjectID,
		PosX:     e.pm.PosX,
		PosY:     e.pm.PosY,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemyAquaman) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyAquaman) getCurrentImagePointer() *int32 {
	n := (e.count / aquamanDelays[e.state])
	if n >= len(e.images[e.state]) {
		n = len(e.images[e.state]) - 1
	}
	return &e.images[e.state][n]
}
