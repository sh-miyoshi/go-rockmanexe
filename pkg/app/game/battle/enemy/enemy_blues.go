package enemy

import (
	"math/rand"

	"github.com/cockroachdb/errors"

	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	deleteanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/delete"
	localanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/local"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/skillcore"
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
	bluesDelays = [bluesActTypeMax]int{1, 2, 4, 4, 4, 4, 4, 1}
)

type enemyBlues struct {
	pm               EnemyParam
	state            int
	count            int
	waitCount        int
	nextState        int
	targetPos        point.Point
	isTargetPosMoved bool
	moveNum          int
	images           [bluesActTypeMax][]int
	edgeAtkCount     int
	atkIDs           []string
	isCharReverse    bool
}

func (e *enemyBlues) Init(objID string) error {
	e.pm.ObjectID = objID
	e.state = bluesActTypeStand
	e.count = 0
	e.waitCount = 20
	e.nextState = bluesActTypeMove
	e.targetPos = emptyPos
	e.moveNum = 2
	e.isTargetPosMoved = false
	e.edgeAtkCount = 0
	e.isCharReverse = false

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

	e.images[bluesActTypeDeltaRayEdge] = make([]int, 1)
	e.images[bluesActTypeDeltaRayEdge][0] = tmp[0]

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
			if !e.targetPos.Equal(emptyPos) {
				if !battlecommon.MoveObjectDirect(
					&e.pm.Pos,
					e.targetPos,
					-1, // プレイヤーのパネルでも移動可能
					true,
					field.GetPanelInfo,
				) {
					// 移動に失敗したら、移動からやり直し
					logger.Debug("Forte move failed. retry")
					return e.clearState()
				}
				e.targetPos = emptyPos
				if e.waitCount == 0 {
					e.waitCount = 20
				}
				return e.stateChange(forteActTypeStand)
			}

			moveRandom(&e.pm.Pos)
			e.waitCount = 40

			e.moveNum--
			if e.moveNum <= 0 {
				if debugFlag {
					e.moveNum = 3
					e.nextState = bluesActTypeDeltaRayEdge
					return e.stateChange(bluesActTypeStand)
				}

				e.moveNum = rand.Intn(3) + 3
				// WIP: 確率によって行動を変える
			}

			return e.stateChange(bluesActTypeStand)
		}
	case bluesActTypeWideSword:
		if e.count == 0 && !e.isTargetPosMoved {
			e.isTargetPosMoved = true

			// Move to attack position
			objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				return e.clearState()
			}
			// TODO: 一旦右上だが、右下も候補にする
			targetPos := point.Point{X: objs[0].Pos.X + 1, Y: objs[0].Pos.Y - 1}
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = bluesActTypeWideSword
				return e.stateChange(bluesActTypeMove)
			}
		}

		if e.count == 1*bluesDelays[bluesActTypeWideSword] {
			logger.Debug("Blues Wide Sword Attack")
			localanim.AnimNew(skill.Get(resources.SkillWideSword, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      forteAtkPower[e.state],
				TargetType: damage.TargetPlayer,
			}))
		}

		if e.count == 6*bluesDelays[bluesActTypeWideSword] {
			return e.clearState()
		}
	case bluesActTypeFighterSword:
		if e.count == 0 && !e.isTargetPosMoved {
			e.isTargetPosMoved = true

			// Move to attack position
			objs := localanim.ObjAnimGetObjs(objanim.Filter{ObjType: objanim.ObjTypePlayer})
			if len(objs) == 0 {
				// エラー処理
				logger.Info("Failed to get player position")
				return e.clearState()
			}
			tx := 0
			for x := 0; x < battlecommon.FieldNum.X; x++ {
				if field.GetPanelInfo(point.Point{X: x, Y: objs[0].Pos.Y}).Type == battlecommon.PanelTypeEnemy {
					tx = x
					break
				}
			}

			targetPos := point.Point{X: tx, Y: objs[0].Pos.Y}
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = bluesActTypeFighterSword
				return e.stateChange(bluesActTypeMove)
			}
		}

		if e.count == 1*bluesDelays[bluesActTypeFighterSword] {
			logger.Debug("Blues Fighter Sword Attack")
			localanim.AnimNew(skill.Get(resources.SkillFighterSword, skillcore.Argument{
				OwnerID:    e.pm.ObjectID,
				Power:      forteAtkPower[e.state],
				TargetType: damage.TargetPlayer,
			}))
		}

		if e.count == 6*bluesDelays[bluesActTypeWideSword] {
			return e.clearState()
		}
	case bluesActTypeDeltaRayEdge:
		if e.count == 0 {
			localanim.AnimNew(effect.Get(resources.EffectTypeSpecialStart, e.pm.Pos, 0))
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
		{X: -5, Y: -20},  // Stand
		{X: -5, Y: -20},  // Move
		{X: -20, Y: -20}, // WideSword
		{X: -20, Y: -20}, // FighterSword
		{X: -20, Y: -20}, // SonicBoom
		{X: -20, Y: -20}, // DeltaRayEdge
		{X: -20, Y: -20}, // BehindSlash
		{X: 0, Y: 0},     // Damage
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

func (e *enemyBlues) clearState() (bool, error) {
	e.waitCount = 20
	e.nextState = forteActTypeMove
	e.moveNum = 3 + rand.Intn(3)
	e.targetPos = emptyPos
	e.isTargetPosMoved = false

	return e.stateChange(forteActTypeStand)
}
