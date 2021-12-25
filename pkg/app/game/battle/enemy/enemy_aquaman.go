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
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
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
	images          [aquamanActTypeMax][]int32
	count           int
	state           int
	nextState       int
	waitCount       int
	invincibleCount int
	actID           string
	waterPipeObjIDs []string
	moveNum         int
	targetPosX      int
	targetPosY      int
	beforeAction    int
}

func (e *enemyAquaman) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = aquamanActTypeStand
	e.waitCount = 60
	e.nextState = aquamanActTypeMove
	e.waterPipeObjIDs = []string{}
	e.moveNum = rand.Intn(2) + 2
	e.targetPosX = -1
	e.targetPosY = -1
	e.beforeAction = -1

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

	fname = name + "_create" + ext
	e.images[aquamanActTypeCreate] = make([]int32, 1)
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
		e.images[i] = []int32{}
	}

	for _, id := range e.waterPipeObjIDs {
		anim.Delete(id)
	}
	e.waterPipeObjIDs = []string{}
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
			if e.targetPosX != -1 && e.targetPosY != -1 {
				if battlecommon.MoveObjectDirect(
					&e.pm.PosX,
					&e.pm.PosY,
					e.targetPosX,
					e.targetPosY,
					field.PanelTypeEnemy,
					true,
					field.GetPanelInfo,
				) {
					e.targetPosX = -1
					e.targetPosY = -1
					e.count = 0
					e.waitCount = 20
					e.state = aquamanActTypeStand
					return false, nil
				}
			}

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
			e.waitCount = 20
			e.state = aquamanActTypeStand
			e.moveNum--
			if e.moveNum <= 0 {
				// Select attack
				n := rand.Intn(100)
				if n < 20 && e.beforeAction != aquamanActTypeCreate {
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
			tx := 1
			ty := 1
			if len(objs) > 0 {
				tx = objs[0].PosX
				ty = objs[0].PosY
			}
			if e.pm.PosX != tx+3 || e.pm.PosY != ty {
				e.targetPosX = tx + 3
				e.targetPosY = ty
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
			if e.pm.PosX == 3 && (e.pm.PosY == 0 || e.pm.PosY == 2) {
				e.targetPosX = 5
				e.targetPosY = 1
				e.state = aquamanActTypeMove
				e.count = 0
				return false, nil
			}
		}

		if e.count == 5 {
			obj := &object.WaterPipe{}
			pm := object.ObjectParam{
				PosX:          3,
				HP:            500,
				OnwerCharType: objanim.ObjTypeEnemy,
				AttackNum:     5,
				Interval:      150,
				Power:         20,
			}
			pm.PosY = 0
			if err := obj.Init(e.pm.ObjectID, pm); err != nil {
				return false, fmt.Errorf("water pipe create failed: %w", err)
			}
			e.waterPipeObjIDs = append(e.waterPipeObjIDs, objanim.New(obj))
			obj = &object.WaterPipe{}
			pm.PosY = 2
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
	x, y := battlecommon.ViewPos(e.pm.PosX, e.pm.PosY)
	img := e.getCurrentImagePointer()

	ofs := [aquamanActTypeMax][]int32{
		{0, 0},    // stand
		{0, 0},    // move
		{-20, 10}, // shot
		{0, 0},    // damage
		{0, 0},    // bomb
		{0, 0},    // create
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
