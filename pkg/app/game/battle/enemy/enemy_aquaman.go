package enemy

import (
	"fmt"
	"math/rand"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
)

const (
	aquamanActTypeStand = iota
	aquamanActTypeMove
	aquamanActTypeShot
	aquamanActTypeDamage
	aquamanActTypeBomb
	aquamanActTypeCreate

	aquamanActTypeMax
)

var (
	aquamanDelays = [aquamanActTypeMax]int{8, 4, 4, 4, 3, 1}
)

type enemyAquaman struct {
	pm              EnemyParam
	images          [aquamanActTypeMax][]int
	count           int
	state           int
	nextState       int
	waitCount       int
	invincibleCount int
	actID           string
	waterPipeObjIDs []string
	moveNum         int
	targetPos       common.Point
	beforeAction    int
}

func (e *enemyAquaman) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = aquamanActTypeStand
	e.waitCount = 60
	e.nextState = aquamanActTypeMove
	e.waterPipeObjIDs = []string{}
	e.moveNum = rand.Intn(2) + 2
	e.targetPos = common.Point{X: -1, Y: -1}
	e.beforeAction = -1

	// Load Images
	name, ext := GetStandImageFile(IDAquaman)

	fname := name + "_stand" + ext
	e.images[aquamanActTypeStand] = make([]int, 9)
	if res := dxlib.LoadDivGraph(fname, 9, 9, 1, 62, 112, e.images[aquamanActTypeStand]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_move" + ext
	e.images[aquamanActTypeMove] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 64, 126, e.images[aquamanActTypeMove]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_shot" + ext
	e.images[aquamanActTypeShot] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 104, 90, e.images[aquamanActTypeShot]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_damage" + ext
	e.images[aquamanActTypeDamage] = make([]int, 1)
	if res := dxlib.LoadDivGraph(fname, 1, 1, 1, 70, 86, e.images[aquamanActTypeDamage]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_bomb" + ext
	e.images[aquamanActTypeBomb] = make([]int, 5)
	if res := dxlib.LoadDivGraph(fname, 5, 5, 1, 100, 124, e.images[aquamanActTypeBomb]); res == -1 {
		return fmt.Errorf("failed to load image: %s", fname)
	}

	fname = name + "_create" + ext
	e.images[aquamanActTypeCreate] = make([]int, 1)
	if res := dxlib.LoadDivGraph(fname, 1, 1, 1, 80, 92, e.images[aquamanActTypeCreate]); res == -1 {
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
		e.images[i] = []int{}
	}

	for _, id := range e.waterPipeObjIDs {
		anim.Delete(id)
	}
	e.waterPipeObjIDs = []string{}
}

func (e *enemyAquaman) Process() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		battlecommon.NewDelete(*img, e.pm.Pos, false)
		anim.New(effect.Get(effect.TypeExplode, e.pm.Pos, 0))
		*img = -1 // DeleteGraph at delete animation
		return true, nil
	}

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
			if e.targetPos.X != -1 && e.targetPos.Y != -1 {
				if battlecommon.MoveObjectDirect(
					&e.pm.Pos,
					e.targetPos,
					field.PanelTypeEnemy,
					true,
					field.GetPanelInfo,
				) {
					e.targetPos = common.Point{X: -1, Y: -1}
					e.count = 0
					e.waitCount = 20
					e.state = aquamanActTypeStand
					return false, nil
				}
			}

			for i := 0; i < 10; i++ {
				next := common.Point{
					X: rand.Intn(field.FieldNum.X/2) + field.FieldNum.X/2,
					Y: rand.Intn(field.FieldNum.Y),
				}
				if battlecommon.MoveObjectDirect(
					&e.pm.Pos,
					next,
					field.PanelTypeEnemy,
					true,
					field.GetPanelInfo,
				) {
					break
				}
			}
			e.waitCount = 20
			e.state = aquamanActTypeStand
			e.moveNum--
			if e.moveNum <= 0 {
				// Select attack
				n := rand.Intn(100)

				if n < 20 && !e.pipeExists() {
					e.waterPipeObjIDs = []string{}
					e.nextState = aquamanActTypeCreate
					e.moveNum = rand.Intn(2) + 2
				} else if n < 50 {
					e.nextState = aquamanActTypeBomb
					e.moveNum = rand.Intn(2) + 2
				} else {
					e.nextState = aquamanActTypeShot
					e.moveNum = rand.Intn(2) + 1
				}
				e.beforeAction = e.nextState
			}

			e.count = 0
			return false, nil
		}
	case aquamanActTypeShot:
		if e.count == 0 {
			// Move to attack position
			objs := objanim.GetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			t := common.Point{X: 1, Y: 1}
			if len(objs) > 0 {
				t = objs[0].Pos
			}
			if e.pm.Pos.X != t.X+3 || e.pm.Pos.Y != t.Y {
				e.targetPos.X = t.X + 3
				e.targetPos.Y = t.Y
				e.state = aquamanActTypeMove
				e.count = 0
				logger.Debug("%+v", e)
				return false, nil
			}
		}

		if e.count == 0 {
			e.actID = anim.New(skill.Get(skill.SkillAquamanShot, skill.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      10,
				TargetType: damage.TargetPlayer,
			}))
		}

		if !anim.IsProcessing(e.actID) {
			e.waitCount = 60
			e.state = aquamanActTypeStand
			e.nextState = aquamanActTypeMove
			e.count = 0
			return false, nil
		}
	case aquamanActTypeDamage:
		if e.count == 4*aquamanDelays[aquamanActTypeDamage] {
			e.waitCount = 20
			e.state = aquamanActTypeStand
			e.nextState = aquamanActTypeMove
			e.count = 0
			return false, nil
		}
	case aquamanActTypeBomb:
		if e.count == 3*aquamanDelays[aquamanActTypeBomb] {
			// ボム登録
			anim.New(skill.Get(skill.SkillWaterBomb, skill.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      50,
				TargetType: damage.TargetPlayer,
			}))
		}

		if e.count == 6*aquamanDelays[aquamanActTypeBomb] {
			e.waitCount = 90
			e.state = aquamanActTypeStand
			e.nextState = aquamanActTypeMove
			e.count = 0
			return false, nil
		}
	case aquamanActTypeCreate:
		if e.count == 0 {
			if e.pm.Pos.X == 3 && (e.pm.Pos.Y == 0 || e.pm.Pos.Y == 2) {
				e.targetPos.X = 5
				e.targetPos.Y = 1
				e.state = aquamanActTypeMove
				e.count = 0
				return false, nil
			}
		}

		if e.count == 5 {
			obj := &object.WaterPipe{}
			pm := object.ObjectParam{
				HP:            500,
				OnwerCharType: objanim.ObjTypeEnemy,
				AttackNum:     5,
				Interval:      150,
				Power:         20,
			}
			pm.Pos.X = 3
			pm.Pos.Y = 0
			if err := obj.Init(e.pm.ObjectID, pm); err != nil {
				return false, fmt.Errorf("water pipe create failed: %w", err)
			}
			e.waterPipeObjIDs = append(e.waterPipeObjIDs, objanim.New(obj))
			obj = &object.WaterPipe{}
			pm.Pos.Y = 2
			if err := obj.Init(e.pm.ObjectID, pm); err != nil {
				return false, fmt.Errorf("water pipe create failed: %w", err)
			}
			e.waterPipeObjIDs = append(e.waterPipeObjIDs, objanim.New(obj))
		}

		if e.count == 50 {
			e.waitCount = 60
			e.state = aquamanActTypeStand
			e.nextState = aquamanActTypeMove
			e.count = 0
		}
	}

	e.count++
	return false, nil
}

func (e *enemyAquaman) Draw() {
	if e.invincibleCount/5%2 != 0 {
		return
	}

	if e.state == aquamanActTypeShot && e.count == 0 {
		// 移動が必要な際にShotの画像を表示したくないため
		return
	}

	// Show Enemy Images
	view := battlecommon.ViewPos(e.pm.Pos)
	img := e.getCurrentImagePointer()

	ofs := [aquamanActTypeMax][]int{
		{0, 0},    // stand
		{0, 0},    // move
		{-20, 10}, // shot
		{0, 0},    // damage
		{0, 0},    // bomb
		{0, 0},    // create
	}

	dxlib.DrawRotaGraph(view.X+ofs[e.state][0], view.Y+ofs[e.state][1], 1, 0, *img, true)

	// Show HP
	if e.pm.HP > 0 {
		draw.Number(view.X, view.Y+40, e.pm.HP, draw.NumberOption{
			Color:    draw.NumberColorWhiteSmall,
			Centered: true,
		})
	}
}

func (e *enemyAquaman) DamageProc(dm *damage.Damage) bool {
	if dm == nil {
		return false
	}

	if e.invincibleCount > 0 {
		return false
	}

	if dm.TargetType&damage.TargetEnemy != 0 {
		e.pm.HP -= dm.Power
		anim.New(effect.Get(dm.HitEffectType, e.pm.Pos, 5))

		for i := 0; i < dm.PushLeft; i++ {
			if !battlecommon.MoveObject(&e.pm.Pos, common.DirectLeft, field.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}
		for i := 0; i < dm.PushRight; i++ {
			if !battlecommon.MoveObject(&e.pm.Pos, common.DirectRight, field.PanelTypeEnemy, true, field.GetPanelInfo) {
				break
			}
		}

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
		Pos:      e.pm.Pos,
		AnimType: anim.AnimTypeObject,
	}
}

func (e *enemyAquaman) GetObjectType() int {
	return objanim.ObjTypeEnemy
}

func (e *enemyAquaman) getCurrentImagePointer() *int {
	n := (e.count / aquamanDelays[e.state])
	if n >= len(e.images[e.state]) {
		n = len(e.images[e.state]) - 1
	}
	return &e.images[e.state][n]
}

func (e *enemyAquaman) pipeExists() bool {
	for _, id := range e.waterPipeObjIDs {
		if objanim.IsProcessing(id) {
			return true
		}
	}
	return false
}
