package enemy

import (
	"math/rand"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/draw"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/manager"
	objanim "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/anim/object"
	battlecommon "github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/common/deleteanim"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/damage"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/effect"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/field"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/object"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/game/battle/skill"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/resources"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/app/system"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/dxlib"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/logger"
	"github.com/sh-miyoshi/go-rockmanexe/pkg/utils/point"
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
	coldmanDelays = [coldmanActTypeMax]int{1, 1, 2, 7, 1, 3, 5}
)

type enemyColdman struct {
	pm           EnemyParam
	images       [coldmanActTypeMax][]int
	count        int
	state        int
	nextState    int
	waitCount    int
	moveNum      int
	cubeIDs      []string
	bressIDs     []string
	targetPos    point.Point
	targetCubeID string
	actCount     int
	animMgr      *manager.Manager
}

func (e *enemyColdman) Init(objID string, animMgr *manager.Manager) error {
	e.pm.ObjectID = objID
	e.state = coldmanActTypeStand
	e.waitCount = 60
	e.nextState = coldmanActTypeMove
	e.moveNum = rand.Intn(2) + 2
	e.cubeIDs = []string{}
	e.bressIDs = []string{}
	e.targetPos = point.Point{X: -1, Y: -1}
	e.actCount = 0
	e.animMgr = animMgr

	// Load Images
	name, ext := GetStandImageFile(IDColdman)

	fname := name + "_all" + ext
	tmp := make([]int, 24)
	if res := dxlib.LoadDivGraph(fname, 24, 6, 4, 136, 115, tmp); res == -1 {
		return errors.Newf("failed to load image: %s", fname)
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

func (e *enemyColdman) Update() (bool, error) {
	if e.pm.HP <= 0 {
		// Delete Animation
		img := e.getCurrentImagePointer()
		deleteanim.New(*img, e.pm.Pos, false, e.animMgr)
		e.animMgr.EffectAnimNew(effect.Get(resources.EffectTypeExplode, e.pm.Pos, 0))
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
	case coldmanActTypeStand:
		e.waitCount--
		if e.waitCount <= 0 {
			return e.stateChange(e.nextState)
		}
	case coldmanActTypeMove:
		if e.count == 2*coldmanDelays[coldmanActTypeMove] {
			if e.targetPos.X != -1 && e.targetPos.Y != -1 {
				if !battlecommon.MoveObjectDirect(
					&e.pm.Pos,
					e.targetPos,
					battlecommon.PanelTypeEnemy,
					true,
					field.GetPanelInfo,
				) {
					// 移動に失敗したら、ランダム移動からやり直し
					e.nextState = coldmanActTypeMove
					e.moveNum = rand.Intn(2) + 2
					e.targetCubeID = ""
				}
				e.targetPos = point.Point{X: -1, Y: -1}
				e.waitCount = 20
				return e.stateChange(coldmanActTypeStand)
			}

			e.moveRandom()

			e.waitCount = 60
			e.state = coldmanActTypeStand
			e.moveNum--
			if e.moveNum <= 0 {
				e.moveNum = rand.Intn(2) + 2

				// Action process
				switch e.actCount {
				case 0:
					e.nextState = coldmanActTypeIceCreate
				case 1, 2:
					_, _, exists := e.findShootCube()
					if exists {
						e.nextState = coldmanActTypeIceShoot
					} else {
						e.nextState = coldmanActTypeBless
					}
				}
				e.actCount = (e.actCount + 1) % 3
			}
		}
	case coldmanActTypeIceCreate:
		if e.count == 0 {
			field.SetBlackoutCount(90)
			skill.SetChipNameDraw("アイスキューブ", false)

			if err := e.createCube(); err != nil {
				return false, err
			}
		}
		if e.count == 60 {
			for _, id := range e.cubeIDs {
				e.animMgr.DeactivateAnim(id)
			}
			e.moveNum = rand.Intn(2) + 2
			e.nextState = coldmanActTypeMove
			return e.stateChange(coldmanActTypeMove)
		}
	case coldmanActTypeIceShoot:
		if e.count == 0 && e.targetCubeID == "" {
			// キューブの前に移動
			targetPos, id, exists := e.findShootCube()
			if !exists {
				// 対象のキューブが見つからない場合、ランダム移動からやり直し
				e.nextState = coldmanActTypeMove
				e.waitCount = 20
				return e.stateChange(coldmanActTypeStand)
			}
			e.targetCubeID = id
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = coldmanActTypeIceShoot
				return e.stateChange(coldmanActTypeMove)
			}
			return false, nil
		}

		if e.count == 3*coldmanDelays[coldmanActTypeIceShoot] {
			e.targetCubeID = ""
			// PUSH
			dm := damage.Damage{
				ID:            uuid.New().String(),
				DamageType:    damage.TypePosition,
				Pos:           point.Point{X: e.pm.Pos.X - 1, Y: e.pm.Pos.Y},
				Power:         1, // debug(0でもいいかも)
				TTL:           1,
				TargetObjType: damage.TargetPlayer | damage.TargetEnemy,
				HitEffectType: resources.EffectTypeNone,
				ShowHitArea:   false,
				StrengthType:  damage.StrengthHigh,
				PushLeft:      battlecommon.FieldNum.X,
				Element:       damage.ElementNone,
			}
			e.animMgr.DamageManager().New(dm)
		}

		if e.count == 6*coldmanDelays[coldmanActTypeIceShoot] {
			e.waitCount = 20
			e.nextState = coldmanActTypeMove
			return e.stateChange(coldmanActTypeStand)
		}
	case coldmanActTypeBodyBlow:
		system.SetError("TODO: not implemented yet")
	case coldmanActTypeBless:
		if e.count == 0 {
			targetPos := point.Point{X: 5, Y: 1}
			if !targetPos.Equal(e.pm.Pos) {
				e.targetPos = targetPos
				e.nextState = coldmanActTypeBless
				return e.stateChange(coldmanActTypeMove)
			}
		}

		if e.count == 2*coldmanDelays[coldmanActTypeBless] {
			e.createBress()

			e.moveNum = rand.Intn(2) + 2
			e.waitCount = 60
			e.nextState = coldmanActTypeMove
			return e.stateChange(coldmanActTypeStand)
		}
	case coldmanActTypeDamage:
		if e.count == 4*coldmanDelays[coldmanActTypeDamage] {
			e.waitCount = 20
			e.nextState = coldmanActTypeMove
			return e.stateChange(coldmanActTypeStand)
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

	ofs := [coldmanActTypeMax]point.Point{
		{X: 0, Y: 0},  // Stand
		{X: 0, Y: 0},  // IceCreate
		{X: 0, Y: 0},  // Move
		{X: 0, Y: 0},  // IceShoot
		{X: 0, Y: 0},  // BodyBlow
		{X: 0, Y: 0},  // Bless
		{X: 20, Y: 0}, // Damage
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

func (e *enemyColdman) DamageProc(dm *damage.Damage) bool {
	if damageProc(dm, &e.pm) {
		if dm.StrengthType == damage.StrengthNone {
			return true
		}

		e.state = coldmanActTypeDamage
		if dm.StrengthType == damage.StrengthHigh {
			e.pm.InvincibleCount = battlecommon.PlayerDefaultInvincibleTime
		}
		e.count = 0
		return true
	}

	return false
}

func (e *enemyColdman) GetParam() objanim.Param {
	return objanim.Param{
		Param: anim.Param{
			ObjID: e.pm.ObjectID,
			Pos:   e.pm.Pos,
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

func (e *enemyColdman) AddBarrier(hp int) {}

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
		next := point.Point{
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
	// 前のアイスキューブがあるなら削除する
	if len(e.cubeIDs) > 0 {
		for _, id := range e.cubeIDs {
			e.animMgr.AnimDelete(id)
		}
		e.cubeIDs = []string{}
	}

	// 特定のパターンで2個生成
	// 原作では3個だが実装の都合上2個にする

	// パターン1
	//   -------------
	//   | x |   |   |
	//   -------------
	//   |   |   |   |
	//   -------------
	//   |   | x |   |
	//   -------------

	// パターン2
	//   -------------
	//   |   |   |   |
	//   -------------
	//   | x |   |   |
	//   -------------
	//   | x |   |   |
	//   -------------
	patterns := [][]point.Point{
		{
			point.Point{X: 3, Y: 0},
			point.Point{X: 4, Y: 2},
		},
		{
			point.Point{X: 3, Y: 1},
			point.Point{X: 3, Y: 2},
		},
	}

	index := rand.Intn(len(patterns))
	ptn := patterns[index]

	for _, pos := range ptn {
		pm := object.ObjectParam{
			Pos:           pos,
			HP:            60,
			OnwerCharType: objanim.ObjTypeEnemy,
		}
		obj := &object.IceCube{}
		if err := obj.Init(e.pm.ObjectID, pm, e.animMgr); err != nil {
			return errors.Wrap(err, "failed to init ice cube")
		}
		id := e.animMgr.ObjAnimNew(obj)
		e.animMgr.SetActiveAnim(id)
		e.cubeIDs = append(e.cubeIDs, id)
	}

	return nil
}

func (e *enemyColdman) createBress() error {
	// 前のブレスがあるなら削除する
	if len(e.bressIDs) > 0 {
		for _, id := range e.bressIDs {
			e.animMgr.AnimDelete(id)
		}
		e.bressIDs = []string{}
	}

	for y := 0; y < battlecommon.FieldNum.Y; y++ {
		pos := point.Point{X: 4, Y: y}
		// もしObjectがあれば生成しない
		if e.animMgr.ObjAnimExistsObject(pos) != "" {
			continue
		}

		pm := object.ObjectParam{
			Pos:           pos,
			HP:            10,
			OnwerCharType: objanim.ObjTypeEnemy,
		}

		obj := &object.ColdBress{}
		if err := obj.Init(e.pm.ObjectID, pm, e.animMgr); err != nil {
			return errors.Wrap(err, "failed to init cold bress")
		}
		id := e.animMgr.ObjAnimNew(obj)
		e.bressIDs = append(e.bressIDs, id)
	}

	return nil
}

func (e *enemyColdman) stateChange(next int) (bool, error) {
	logger.Info("change coldman state to %d", next)
	e.state = next
	e.count = 0

	return false, nil
}

func (e *enemyColdman) findShootCube() (point.Point, string, bool) {
	playerPos := e.animMgr.ObjAnimGetObjPos(e.pm.PlayerID)
	for _, id := range e.cubeIDs {
		pos := e.animMgr.ObjAnimGetObjPos(id)
		if pos.Y == playerPos.Y {
			targetPos := point.Point{X: pos.X + 1, Y: pos.Y}
			return targetPos, id, true
		}
	}
	return point.Point{X: -1, Y: -1}, "", false
}
